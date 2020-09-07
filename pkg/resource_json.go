package tutorial

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

type BasicValue struct {
	Score int
}

var client = NewDriver("json")

func resourceTutorialJson() *schema.Resource {
	return &schema.Resource{
		Create: resourceTutorialJsonCreate,
		Read:   resourceTutorialJsonRead,
		Update: resourceTutorialJsonUpdate,
		Delete: resourceTutorialJsonDelete,
		Schema: map[string]*schema.Schema{
			"score": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceTutorialJsonCreate(d *schema.ResourceData, meta interface{}) error {
	score := d.Get("score").(int)
	value := &BasicValue{
		Score: score,
	}

	id, err := client.Create(value)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceTutorialJsonRead(d, meta)
}

func resourceTutorialJsonRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	data, err := client.Read(id)
	if err != nil {
		if IsStateNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	value := &BasicValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	d.Set("score", value.Score) // nolint

	return nil
}

func resourceTutorialJsonUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	data, err := client.Read(id)
	if err != nil {
		if IsStateNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	value := &BasicValue{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if d.HasChange("score") {
		value.Score = d.Get("score").(int)
	}

	if err := client.Update(id, value); err != nil {
		return err
	}

	return resourceTutorialJsonRead(d, meta)
}

func resourceTutorialJsonDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	if err := client.Delete(id); err != nil {
		if !IsStateNotFoundError(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}

const emptyID = ""

func init() {
	rand.Seed(time.Now().UnixNano())
}

type NotFoundError error

func IsStateNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(NotFoundError)
	return ok
}

func NewDriver(typeName string) *Driver {
	return &Driver{typeName: typeName}
}

type Driver struct {
	typeName string
}

func (d *Driver) Create(value interface{}) (string, error) {
	id := d.generateID()
	statePath := d.stateFilePath(id)

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return emptyID, fmt.Errorf("marshaling to JSON is failed: %s", err)
	}

	err = ioutil.WriteFile(statePath, data, 0644)
	if err != nil {
		return emptyID, fmt.Errorf("writing to file is failed: %s", err)
	}

	return id, nil
}

func (d *Driver) Read(id string) ([]byte, error) {
	if !d.stateExists(id) {
		return nil, NotFoundError(errors.New("state not found"))
	}
	statePath := d.stateFilePath(id)
	return ioutil.ReadFile(statePath)
}

func (d *Driver) Update(id string, value interface{}) error {
	if !d.stateExists(id) {
		return NotFoundError(errors.New("state not found"))
	}
	statePath := d.stateFilePath(id)

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling to JSON is failed: %s", err)
	}

	err = ioutil.WriteFile(statePath, data, 0644)
	if err != nil {
		return fmt.Errorf("writing to file is failed: %s", err)
	}

	return nil
}

func (d *Driver) Delete(id string) error {
	if !d.stateExists(id) {
		return NotFoundError(errors.New("state not found"))
	}
	statePath := d.stateFilePath(id)

	if err := os.Remove(statePath); err != nil {
		return fmt.Errorf("removing file is failed: %s", err)
	}
	return nil
}

func (d *Driver) generateID() string {
	return fmt.Sprintf("%d", rand.Int31())
}

func (d *Driver) stateFilePath(id string) string {
	return fmt.Sprintf("%s-%s.json", d.typeName, id)
}

func (d *Driver) stateExists(id string) bool {
	exists := true
	if _, err := os.Stat(d.stateFilePath(id)); os.IsNotExist(err) {
		exists = false
	}
	return exists
}
