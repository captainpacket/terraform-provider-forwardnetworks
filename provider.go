package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-forwardnetworks/forwardnetworks"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FORWARDNETWORKS_USERNAME", nil),
				Description: "The username for the Forward Networks API.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("FORWARDNETWORKS_PASSWORD", nil),
				Description: "The password for the Forward Networks API.",
			},
			"apphost": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FORWARDNETWORKS_APPHOST", nil),
				Description: "The base URL for the Forward Networks API.",
				Default:	 "fwd.app",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set this to true to allow insecure connections to the Forward Networks API.",
			},
		},
		 ResourcesMap: map[string]*schema.Resource{
         	"forwardnetworks_proxy": resourceProxy(),
			"forwardnetworks_cloud": resourceCloudAccount(),
			"forwardnetworks_network": resourceNetworks(),
			"forwardnetworks_collector": resourceCollectors(),
			"forwardnetworks_collection_schedule": resourceCollectionSchedule(),
			"forwardnetworks_collection": resourceCollection(),
			"forwardnetworks_check": resourceCheck(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"forwardnetworks_version": dataSourceVersion(),
			"forwardnetworks_externalid":  dataSourceExternalId(),
			"forwardnetworks_cloud":  dataSourceCloudAccounts(),
			"forwardnetworks_proxy": dataSourceProxy(),
			"forwardnetworks_collection":      dataCollection(),
			"forwardnetworks_snapshot":      dataSourceSnapshots(),
			"forwardnetworks_nqe":             dataSourceForwardNetworksNQEQuery(),
			"forwardnetworks_nqe_execute":  dataSourceForwardNetworksNQEQueryExecution(),
			"forwardnetworks_checks": dataSourceChecks(),
			"forwardnetworks_aws_policy": dataSourceAWSPolicy(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	baseURL := d.Get("apphost").(string)
	insecure := d.Get("insecure").(bool)

	client := forwardnetworks.NewForwardNetworksClient("https://"+baseURL, username, password, insecure)


	var diags diag.Diagnostics

	if insecure {
		warning := "You have enabled insecure mode. This may expose your connection to security risks. Use with caution."
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Insecure mode enabled",
			Detail:   warning,
		})
	}

	return client, diags
}

