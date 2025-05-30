package ovh

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

func resourceCloudProjectKubeNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectKubeNodePoolCreate,
		Read:   resourceCloudProjectKubeNodePoolRead,
		Delete: resourceCloudProjectKubeNodePoolDelete,
		Update: resourceCloudProjectKubeNodePoolUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectKubeNodePoolImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(time.Hour),
			Update:  schema.DefaultTimeout(time.Hour),
			Delete:  schema.DefaultTimeout(time.Hour),
			Read:    schema.DefaultTimeout(5 * time.Minute),
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"kube_id": {
				Type:        schema.TypeString,
				Description: "Kube ID",
				Required:    true,
				ForceNew:    true,
			},
			"autoscale": {
				Type:        schema.TypeBool,
				Description: "Enable auto-scaling for the pool",
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
			},
			"autoscaling_scale_down_unneeded_time_seconds": {
				Description: "scaleDownUnneededTimeSeconds for autoscaling",
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
			},
			"autoscaling_scale_down_unready_time_seconds": {
				Description: "scaleDownUnreadyTimeSeconds for autoscaling",
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
			},
			"autoscaling_scale_down_utilization_threshold": {
				Description: "scaleDownUtilizationThreshold for autoscaling",
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeFloat,
			},
			"anti_affinity": {
				Type:        schema.TypeBool,
				Description: "Enable anti affinity groups for nodes in the pool",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"flavor_name": {
				Type:        schema.TypeString,
				Description: "Flavor name",
				Required:    true,
				ForceNew:    true,
			},
			"desired_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Optional:    true,
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "NodePool resource name",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"max_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Computed:    true,
				Optional:    true,
			},
			"min_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Computed:    true,
				Optional:    true,
			},
			"monthly_billed": {
				Type:        schema.TypeBool,
				Description: "Enable monthly billing on all nodes in the pool",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},

			// computed
			"available_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes which are actually ready in the pool",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Creation date",
				Computed:    true,
			},
			"current_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes present in the pool",
				Computed:    true,
			},
			"flavor": {
				Type:        schema.TypeString,
				Description: "Flavor name",
				Computed:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "Project id",
				Computed:    true,
			},
			"size_status": {
				Type:        schema.TypeString,
				Description: "Status describing the state between number of nodes wanted and available ones",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status",
				Computed:    true,
			},
			"up_to_date_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes with latest version installed in the pool",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Last update date",
				Computed:    true,
			},
			"template": {
				Description: "Node pool template",
				Optional:    true,
				Type:        schema.TypeSet,
				MaxItems:    1,
				Set:         CustomSchemaSetFunc(),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata": {
							Description: "metadata",
							Required:    true,
							Type:        schema.TypeSet,
							MaxItems:    1,
							Set:         CustomSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"finalizers": {
										Description: "finalizers",
										Required:    true,
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"labels": {
										Description: "labels",
										Required:    true,
										Type:        schema.TypeMap,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
									},
									"annotations": {
										Description: "annotations",
										Required:    true,
										Type:        schema.TypeMap,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
									},
								},
							},
						},
						"spec": {
							Description: "spec",
							Required:    true,
							Type:        schema.TypeSet,
							MaxItems:    1,
							Set:         CustomSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unschedulable": {
										Description: "unschedulable",
										Required:    true,
										Type:        schema.TypeBool,
									},
									"taints": {
										Description: "taints",
										Required:    true,
										Type:        schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeMap,
											Set:  schema.HashString,
											ValidateFunc: func(taintInterface interface{}, path string) (warning []string, errorList []error) {
												taint := taintInterface.(map[string]interface{})

												if taint["key"] == nil {
													return nil, []error{errors.New(fmt.Sprintf("key attribute is mandatory for taint: %s", path))}
												}

												if taint["effect"] == nil {
													return nil, []error{errors.New(fmt.Sprintf("effect attribute is mandatory for taint: %s", path))}
												}

												effectString := taint["effect"].(string)
												effect := TaintEffecTypeToID[effectString]
												if effect == NotATaint {
													return nil, []error{fmt.Errorf("effect: %s is not a allowable taint %#v", effectString, TaintEffecTypeToID)}
												}

												return nil, nil
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCloudProjectKubeNodePoolImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("import Id is not service_name/kubeid/poolid formatted")
	}
	serviceName := splitId[0]
	kubeId := splitId[1]
	id := splitId[2]
	d.SetId(id)
	d.Set("kube_id", kubeId)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectKubeNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", serviceName, kubeId)
	params, err := (&CloudProjectKubeNodePoolCreateOpts{}).FromResource(d)
	if err != nil {
		return err
	}
	res := &CloudProjectKubeNodePoolResponse{}

	log.Printf("[DEBUG] Will create nodepool: %+v", params)
	err = config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %w", endpoint, params, err)
	}

	// This is a fix for a weird bug where the nodepool is not immediately available on API
	log.Printf("[DEBUG] Waiting for nodepool %s to be available", res.Id)
	endpoint = fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, res.Id)
	err = helpers.WaitAvailable(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Waiting for nodepool %s to be READY or ERROR", res.Id)
	err = waitForCloudProjectKubeNodePoolWithStateTarget(config.OVHClient, serviceName, kubeId, res.Id, d.Timeout(schema.TimeoutCreate), []string{"READY", "ERROR"})
	if err != nil {
		return fmt.Errorf("timeout while waiting nodepool %s to be READY: %w", res.Id, err)
	}
	log.Printf("[DEBUG] nodepool %s is READY", res.Id)

	d.SetId(res.Id)

	return resourceCloudProjectKubeNodePoolRead(d, meta)
}

func resourceCloudProjectKubeNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, d.Id())
	res := &CloudProjectKubeNodePoolResponse{}

	log.Printf("[DEBUG] Will read nodepool %s from cluster %s in project %s", d.Id(), kubeId, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read nodepool: %+v", res)
	return nil
}

func resourceCloudProjectKubeNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, d.Id())
	params, err := (&CloudProjectKubeNodePoolUpdateOpts{}).FromResource(d)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Will update nodepool: %#v", *params)
	err = config.OVHClient.Put(endpoint, params, nil)
	if err != nil {
		return fmt.Errorf("calling Put %s with params %v:\n\t %w", endpoint, *params, err)
	}

	log.Printf("[DEBUG] Waiting for nodepool %s to be READY", d.Id())
	err = waitForCloudProjectKubeNodePoolWithStateTarget(config.OVHClient, serviceName, kubeId, d.Id(), d.Timeout(schema.TimeoutUpdate), []string{"READY"})
	if err != nil {
		return fmt.Errorf("timeout while waiting nodepool %s to be READY: %w", d.Id(), err)
	}
	log.Printf("[DEBUG] nodepool %s is READY", d.Id())

	return resourceCloudProjectKubeNodePoolRead(d, meta)
}

func resourceCloudProjectKubeNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, d.Id())

	log.Printf("[DEBUG] Will delete nodepool %s from cluster %s in project %s", d.Id(), kubeId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG] Waiting for nodepool %s to be DELETED", d.Id())
	err = waitForCloudProjectKubeNodePoolDeleted(config.OVHClient, serviceName, kubeId, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("timeout while waiting nodepool %s to be DELETED: %v", d.Id(), err)
	}
	log.Printf("[DEBUG] nodepool %s is DELETED", d.Id())

	d.SetId("")

	return nil
}

func cloudProjectKubeNodePoolExists(serviceName, kubeId, id string, client *ovhwrap.Client) error {
	res := &CloudProjectKubeNodePoolResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, id)
	return client.Get(endpoint, res)
}

func waitForCloudProjectKubeNodePoolWithStateTarget(client *ovhwrap.Client, serviceName, kubeId, id string, timeout time.Duration, stateTargets []string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INSTALLING", "UPDATING", "REDEPLOYING", "RESIZING", "DOWNSCALING", "UPSCALING", "UNKNOWN"},
		Target:  stateTargets,
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectKubeNodePoolResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, id)
			err := client.Get(endpoint, res)
			if err != nil {
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectKubeNodePoolDeleted(client *ovhwrap.Client, serviceName, kubeId, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectKubeNodePoolResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, id)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				} else {
					return res, "", err
				}
			}

			return res, res.Status, nil
		},
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}
