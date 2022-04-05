// Copyright 2021. Clumio, Inc.

// Package clumio_protection_group contains the resource definition and CRUD implementation.
package clumio_protection_group

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/clumio-code/clumio-go-sdk/config"
	protectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	prefixFilter = &schema.Resource{
		Schema: map[string]*schema.Schema{
			schemaExcludedSubPrefixes: {
				Type:        schema.TypeSet,
				Description: "List of subprefixes to exclude from the prefix.",
				Set:         common.SchemaSetHashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			schemaPrefix: {
				Type:        schema.TypeString,
				Description: "Prefix to include.",
				Optional:    true,
			},
		},
	}
	objectFilter = &schema.Resource{
		Schema: map[string]*schema.Schema{
			schemaLatestVersionOnly: {
				Type:        schema.TypeBool,
				Description: "Whether to back up only the latest object version.",
				Optional:    true,
				Default:     false,
			},
			schemaPrefixFilters: {
				Type:        schema.TypeSet,
				Description: "Prefix Filters.",
				Set:         schema.HashResource(prefixFilter),
				Elem:        prefixFilter,
				Optional:    true,
			},
			schemaStorageClasses: {
				Type: schema.TypeSet,
				Description: "Storage class to include in the backup. If not specified," +
					" then all objects across all storage classes will be backed up." +
					" Valid values are: S3 Standard, S3 Standard-IA," +
					" S3 Intelligent-Tiering, S3 One Zone-IA, S3 Glacier and" +
					" S3 Glacier Deep Archive.",
				Set: common.SchemaSetHashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
)

// ClumioS3ProtectionGroup returns the resource for S3 Protection Group.
func ClumioS3ProtectionGroup() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio S3 Protection Group Resource used to create and manage Protection Groups.",

		CreateContext: clumioProtectionGroupCreate,
		ReadContext:   clumioProtectionGroupRead,
		UpdateContext: clumioProtectionGroupUpdate,
		DeleteContext: clumioProtectionGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			schemaBucketRule: {
				Type: schema.TypeString,
				Description: "Describes the possible conditions for a bucket to be " +
					"automatically added to a protection group. For example: " +
					"{\"aws_tag\":{\"$eq\":{\"key\":\"Environment\", \"value\":\"Prod\"}}}",
				Optional: true,
			},
			schemaDescription: {
				Type:        schema.TypeString,
				Description: "The user-assigned description of the protection group.",
				Optional:    true,
			},
			schemaName: {
				Type:        schema.TypeString,
				Description: "The user-assigned name of the protection group.",
				Required:    true,
			},
			schemaObjectFilter: {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Set:      schema.HashResource(objectFilter),
				Elem:     objectFilter,
				Required: true,
			},
			schemaId: {
				Description: "Protection Group Id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaOrganizationalUnitId: {
				Type: schema.TypeString,
				Description: "The Clumio-assigned ID of the organizational unit" +
					" associated with the protection group.",
				Optional: true,
				Computed: true,
			},
			schemaProtectionStatus: {
				Type: schema.TypeString,
				Description: "The protection status of the protection group. Possible" +
					" values include \"protected\", \"unprotected\", and" +
					" \"unsupported\". If the protection group does not support backups," +
					" then this field has a value of unsupported.",
				Computed: true,
			},
			schemaProtectionInfo: {
				Type:        schema.TypeList,
				Description: "The protection policy applied to this resource.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						schemaInheritingEntityId: {
							Description: "The ID of the entity from which protection" +
								" was inherited.",
							Type:     schema.TypeString,
							Computed: true,
						},
						schemaInheritingEntityType: {
							Description: "The type of the entity from which" +
								" protection was inherited.",
							Type:     schema.TypeString,
							Computed: true,
						},
						schemaPolicyId: {
							Type:        schema.TypeString,
							Description: "ID of policy to apply on Protection Group",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// clumioProtectionGroupCreate handles the Create action for the Clumio Protection
// Group Resource.
func clumioProtectionGroupCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	clumioConfig := common.GetClumioConfigForAPI(client, d)
	protectionGroup := protectionGroups.NewProtectionGroupsV1(clumioConfig)
	name := common.GetStringValue(d, schemaName)
	description := common.GetStringValue(d, schemaDescription)
	bucketRule := common.GetStringValue(d, schemaBucketRule)
	objectFilterIface, ok := d.GetOk(schemaObjectFilter)
	if !ok {
		return diag.Errorf("Required attribute object_filter is not present.")
	}
	objectFilter := mapSchemaObjectFilterToClumioObjectFilter(objectFilterIface)
	response, apiErr := protectionGroup.CreateProtectionGroup(
		models.CreateProtectionGroupV1Request{
			BucketRule:   &bucketRule,
			Description:  &description,
			Name:         &name,
			ObjectFilter: objectFilter,
		})
	if apiErr != nil {
		return diag.Errorf("Error creating Protection Group %v. Error: %v",
			name, string(apiErr.Response))
	}
	err := pollForProtectionGroup(ctx, *response.Id, client.ClumioConfig)
	if err != nil {
		return diag.Errorf(
			"Error reading the created Protection Group: %v. Error: %v", response.Id, err)
	}
	d.SetId(*response.Id)
	return clumioProtectionGroupRead(ctx, d, meta)
}

// clumioProtectionGroupUpdate handles the Update action for the Clumio Protection
// Group Resource.
func clumioProtectionGroupUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	clumioConfig := common.GetClumioConfigForAPI(client, d)
	protectionGroup := protectionGroups.NewProtectionGroupsV1(clumioConfig)
	name := common.GetStringValue(d, schemaName)
	description := common.GetStringValue(d, schemaDescription)
	bucketRule := common.GetStringValue(d, schemaBucketRule)
	objectFilterIface, ok := d.GetOk(schemaObjectFilter)
	if !ok {
		return diag.Errorf("Required attribute object_filter is not present.")
	}
	objectFilter := mapSchemaObjectFilterToClumioObjectFilter(objectFilterIface)
	response, apiErr := protectionGroup.UpdateProtectionGroup(d.Id(),
		&models.UpdateProtectionGroupV1Request{
			BucketRule:   &bucketRule,
			Description:  &description,
			Name:         &name,
			ObjectFilter: objectFilter,
		})
	if apiErr != nil {
		return diag.Errorf("Error updating Protection Group %v. Error: %v",
			name, string(apiErr.Response))
	}
	err := pollForProtectionGroup(ctx, *response.Id, client.ClumioConfig)
	if err != nil {
		return diag.Errorf(
			"Error reading the updated Protection Group: %v. Error: %v",
			response.Id, err)
	}
	return clumioProtectionGroupRead(ctx, d, meta)
}

// clumioProtectionGroupCreate handles the Read action for the Clumio Protection
// Group Resource.
func clumioProtectionGroupRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	clumioConfig := common.GetClumioConfigForAPI(client, d)
	protectionGroup := protectionGroups.NewProtectionGroupsV1(clumioConfig)
	readResponse, apiErr := protectionGroup.ReadProtectionGroup(d.Id())
	if apiErr != nil {
		return diag.Errorf("Error creating Protection Group %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	err := d.Set(schemaDescription, *readResponse.Description)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaDescription, err)
	}
	err = d.Set(schemaBucketRule, *readResponse.BucketRule)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaBucketRule, err)
	}
	err = d.Set(schemaName, *readResponse.Name)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaName, err)
	}
	schemaObjFilter := mapClumioObjectFilterToSchemaObjectFilter(
		readResponse.ObjectFilter)
	err = d.Set(schemaObjectFilter, schemaObjFilter)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaObjectFilter, err)
	}
	err = d.Set(schemaOrganizationalUnitId, *readResponse.OrganizationalUnitId)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaOrganizationalUnitId, err)
	}
	err = d.Set(schemaProtectionStatus, *readResponse.ProtectionStatus)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaProtectionStatus, err)
	}
	var schemaProtectInfo []interface{}
	if readResponse.ProtectionInfo != nil {
		schemaProtectInfo = make([]interface{}, 0)
		protectInfoAttrMap := make(map[string]interface{})
		if readResponse.ProtectionInfo.InheritingEntityId != nil {
			protectInfoAttrMap[schemaInheritingEntityId] =
				*readResponse.ProtectionInfo.InheritingEntityId
		}
		if readResponse.ProtectionInfo.InheritingEntityType != nil {
			protectInfoAttrMap[schemaInheritingEntityType] =
				*readResponse.ProtectionInfo.InheritingEntityType
		}
		if readResponse.ProtectionInfo.PolicyId != nil {
			protectInfoAttrMap[schemaPolicyId] =
				*readResponse.ProtectionInfo.PolicyId
		}
		schemaProtectInfo = append(schemaProtectInfo, protectInfoAttrMap)
	}
	err = d.Set(schemaProtectionInfo, schemaProtectInfo)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaProtectionInfo, err)
	}
	return nil
}

// clumioProtectionGroupCreate handles the Delete action for the Clumio Protection
// Group Resource.
func clumioProtectionGroupDelete(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	clumioConfig := common.GetClumioConfigForAPI(client, d)
	protectionGroup := protectionGroups.NewProtectionGroupsV1(clumioConfig)
	_, apiErr := protectionGroup.DeleteProtectionGroup(d.Id())
	if apiErr != nil {
		return diag.Errorf(
			"Error deleting Protection Group: %v. Error: %v", d.Id(), apiErr)
	}
	return nil
}

// mapClumioObjectFilterToSchemaObjectFilter converts the schema object_filter
// to the model Object Filter
func mapSchemaObjectFilterToClumioObjectFilter(objectFilterIface interface{}) *models.ObjectFilter {
	objectFilters := objectFilterIface.(*schema.Set).List()
	objectFilterAttrMap := objectFilters[0].(map[string]interface{})
	latestVersionOnly := objectFilterAttrMap[schemaLatestVersionOnly].(bool)
	storageClassesIfaceList := objectFilterAttrMap[schemaStorageClasses].(*schema.Set).List()
	storageClasses := make([]*string, 0)
	for _, storageClassIface := range storageClassesIfaceList {
		storageClass := storageClassIface.(string)
		storageClasses = append(storageClasses, &storageClass)
	}
	prefixFilters := objectFilterAttrMap[schemaPrefixFilters].(*schema.Set).List()
	modelPrefixFilters := make([]*models.PrefixFilter, 0)
	for _, prefixFilterIface := range prefixFilters {
		var excludedSubPrefixes []*string
		prefixFilterAttrMap := prefixFilterIface.(map[string]interface{})
		excludedSubPrefixesIface, ok := prefixFilterAttrMap[schemaExcludedSubPrefixes]
		if ok {
			excludedSubPrefixesIfaceList := excludedSubPrefixesIface.(*schema.Set).List()
			excludedSubPrefixes = make([]*string, 0)
			for _, excludedSubPrefixIface := range excludedSubPrefixesIfaceList {
				excludedSubPrefix := excludedSubPrefixIface.(string)
				excludedSubPrefixes = append(excludedSubPrefixes, &excludedSubPrefix)
			}
		}
		var prefix string
		prefixIface, ok := prefixFilterAttrMap[schemaPrefix]
		if ok {
			prefix = prefixIface.(string)
		}
		modelPrefixFilter := &models.PrefixFilter{
			ExcludedSubPrefixes: excludedSubPrefixes,
			Prefix:              &prefix,
		}
		modelPrefixFilters = append(modelPrefixFilters, modelPrefixFilter)
	}
	return &models.ObjectFilter{
		LatestVersionOnly: &latestVersionOnly,
		PrefixFilters:     modelPrefixFilters,
		StorageClasses:    storageClasses,
	}
}

// pollForProtectionGroup polls till the protection group becomes available after create
// or update protection group as they are asynchronous operations.
func pollForProtectionGroup(ctx context.Context, id string, config config.Config) error {
	protectionGroup := protectionGroups.NewProtectionGroupsV1(config)
	interval := time.Duration(intervalInSec) * time.Second
	ticker := time.NewTicker(interval)
	timeout := time.After(time.Duration(timeoutInSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return errors.New("context done")
		case <-ticker.C:
			_, err := protectionGroup.ReadProtectionGroup(id)
			if err != nil {
				errResponseStr := string(err.Response)
				if !strings.Contains(errResponseStr,
					"A resource with the requested ID could not be found.") {
					return errors.New(
						"error reading protection-group which was created")
				}
				continue
			}
			return nil
		case <-timeout:
			return errors.New("polling timeout")
		}
	}
}

// mapClumioObjectFilterToSchemaObjectFilter converts the Object Filter from the
// API to the schema object_filter
func mapClumioObjectFilterToSchemaObjectFilter(
	modelObjectFilter *models.ObjectFilter) *schema.Set {
	schemaObjFilter := &schema.Set{F: schema.HashResource(objectFilter)}
	objFilterAttrMap := make(map[string]interface{})
	if modelObjectFilter.LatestVersionOnly != nil {
		objFilterAttrMap[schemaLatestVersionOnly] = *modelObjectFilter.LatestVersionOnly
	}
	if modelObjectFilter.PrefixFilters != nil {
		prefixFilters := &schema.Set{F: schema.HashResource(prefixFilter)}
		for _, modelPrefixFilter := range modelObjectFilter.PrefixFilters {
			prefixFilterAttrMap := make(map[string]interface{})
			if modelPrefixFilter.Prefix != nil {
				prefixFilterAttrMap[schemaPrefix] = *modelPrefixFilter.Prefix
			}
			if modelPrefixFilter.ExcludedSubPrefixes != nil {
				excludedSubPrefixes := &schema.Set{F: common.SchemaSetHashString}
				for _, excludeSubPrefix := range modelPrefixFilter.ExcludedSubPrefixes {
					excludedSubPrefixes.Add(*excludeSubPrefix)
				}
				prefixFilterAttrMap[schemaExcludedSubPrefixes] = excludedSubPrefixes
			}
			prefixFilters.Add(prefixFilterAttrMap)
		}
		objFilterAttrMap[schemaPrefixFilters] = prefixFilters
	}
	storageClasses := &schema.Set{F: common.SchemaSetHashString}
	for _, storageClass := range modelObjectFilter.StorageClasses {
		storageClasses.Add(*storageClass)
	}
	objFilterAttrMap[schemaStorageClasses] = storageClasses
	schemaObjFilter.Add(objFilterAttrMap)
	return schemaObjFilter
}
