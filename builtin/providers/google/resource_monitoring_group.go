package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/monitoring/v3"
)

func resourceMonitoringGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitoringGroupCreate,
		Read:   resourceMonitoringGroupRead,
		Delete: resourceMonitoringGroupDelete,
		Update: resourceMonitoringGroupUpdate,

		Schema: map[string]*schema.Schema{

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"displayName": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"parentName": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},

			"filter": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"isCluster": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func getGroup(d *schema.ResourceData, config *Config) (*monitoring.Group, error) {
	call := config.clientMonitoring.Projects.Groups.Get(d.Id())
	group, err := call.Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			log.Printf("[WARN] Removing Group %q because it's gone", d.Get("name").(string))
			// The resource doesn't exist anymore
			id := d.Id()
			d.SetId("")

			return nil, fmt.Errorf("Resource %s no longer exists", id)
		}

		return nil, fmt.Errorf("Error reading instance: %s", err)
	}

	return group, nil
}

func resourceMonitoringGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Build the group
	group := &monitoring.Group{}

	if v, ok := d.GetOk("displayName"); ok {
		group.DisplayName = v.(string)
	}

	if v, ok := d.GetOk("parentName"); ok {
		group.ParentName = v.(string)
	}

	if v, ok := d.GetOk("filter"); ok {
		group.Filter = v.(string)
	}

	if v, ok := d.GetOk("isCluster"); ok {
		group.IsCluster = v.(bool)
	}

	// Add the group
	name := fmt.Sprintf("projects/%s", project)
	call := config.clientMonitoring.Projects.Groups.Create(name, group)
	res, err := call.Do()
	if err != nil {
		return err
	}

	d.SetId(res.Name)

	return nil
}

func resourceMonitoringGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	call := config.clientMonitoring.Projects.Groups.Get(d.Id())
	_, err := call.Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't exist anymore
			log.Printf("[WARN] Removing monitoring group %q because it's gone", d.Get("name").(string))
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error reading monitoring group: %s", err)
	}

	return nil
}

func resourceMonitoringGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	group, err := getGroup(d, config)

	call := config.clientMonitoring.Projects.Groups.Update(group.Name, group)
	_, err = call.Do()
	if err != nil {
		return err
	}

	return nil
}

func resourceMonitoringGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Delete the image
	log.Printf("[DEBUG] monitoring group delete request")
	_, err = config.clientCompute.Images.Delete(
		project, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting image: %s", err)
	}

	d.SetId("")
	return nil
}
