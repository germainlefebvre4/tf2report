package diagram

import "strings"

const familyGeneric = "generic"

// serviceFamilies maps a resource type to its service family id. It is a
// read-only table; the first one hundred or so entries cover the common types.
var serviceFamilies = map[string]string{
	// AWS
	"aws_instance": "aws-ec2", "aws_launch_template": "aws-ec2", "aws_autoscaling_group": "aws-autoscaling",
	"aws_vpc": "aws-vpc", "aws_subnet": "aws-subnet", "aws_route_table": "aws-route-table",
	"aws_internet_gateway": "aws-igw", "aws_nat_gateway": "aws-nat", "aws_security_group": "aws-sg",
	"aws_s3_bucket": "aws-s3", "aws_ebs_volume": "aws-ebs", "aws_efs_file_system": "aws-efs",
	"aws_db_instance": "aws-rds", "aws_dynamodb_table": "aws-dynamodb", "aws_elasticache_cluster": "aws-elasticache",
	"aws_lambda_function": "aws-lambda", "aws_lb": "aws-elb", "aws_alb": "aws-elb", "aws_elb": "aws-elb",
	"aws_iam_role": "aws-iam", "aws_iam_policy": "aws-iam", "aws_kms_key": "aws-kms",
	"aws_eks_cluster": "aws-eks", "aws_ecs_cluster": "aws-ecs", "aws_ecr_repository": "aws-ecr",
	"aws_sqs_queue": "aws-sqs", "aws_sns_topic": "aws-sns", "aws_cloudwatch_log_group": "aws-cloudwatch",
	"aws_route53_zone": "aws-route53", "aws_cloudfront_distribution": "aws-cloudfront",
	"aws_api_gateway_rest_api": "aws-apigw",
	// Azure
	"azurerm_linux_virtual_machine": "azure-vm", "azurerm_windows_virtual_machine": "azure-vm",
	"azurerm_virtual_machine_scale_set": "azure-vmss", "azurerm_resource_group": "azure-rg",
	"azurerm_virtual_network": "azure-vnet", "azurerm_subnet": "azure-subnet",
	"azurerm_network_security_group": "azure-nsg", "azurerm_network_interface": "azure-nic",
	"azurerm_public_ip": "azure-pip", "azurerm_lb": "azure-lb", "azurerm_storage_account": "azure-storage",
	"azurerm_storage_container": "azure-storage", "azurerm_key_vault": "azure-keyvault",
	"azurerm_kubernetes_cluster": "azure-aks", "azurerm_container_registry": "azure-acr",
	"azurerm_mssql_server": "azure-sql", "azurerm_mssql_database": "azure-sql",
	"azurerm_postgresql_flexible_server": "azure-postgres", "azurerm_cosmosdb_account": "azure-cosmos",
	"azurerm_redis_cache": "azure-redis", "azurerm_linux_function_app": "azure-functions",
	"azurerm_service_plan": "azure-appservice", "azurerm_linux_web_app": "azure-appservice",
	"azurerm_dns_zone": "azure-dns", "azurerm_log_analytics_workspace": "azure-monitor",
	"azurerm_user_assigned_identity": "azure-identity", "azurerm_servicebus_namespace": "azure-servicebus",
	"azurerm_application_gateway": "azure-appgw",
	// GCP
	"google_compute_instance": "gcp-gce", "google_compute_network": "gcp-vpc",
	"google_compute_subnetwork": "gcp-subnet", "google_compute_firewall": "gcp-firewall",
	"google_compute_address": "gcp-address", "google_compute_disk": "gcp-disk",
	"google_compute_forwarding_rule": "gcp-lb", "google_compute_instance_group": "gcp-gce",
	"google_storage_bucket": "gcp-gcs", "google_sql_database_instance": "gcp-sql",
	"google_container_cluster": "gcp-gke", "google_container_node_pool": "gcp-gke",
	"google_cloud_run_service": "gcp-run", "google_cloudfunctions_function": "gcp-functions",
	"google_pubsub_topic": "gcp-pubsub", "google_bigquery_dataset": "gcp-bigquery",
	"google_kms_key_ring": "gcp-kms", "google_project_iam_member": "gcp-iam",
	"google_service_account": "gcp-iam", "google_dns_managed_zone": "gcp-dns",
	"google_redis_instance": "gcp-memorystore", "google_secret_manager_secret": "gcp-secret",
	"google_artifact_registry_repository": "gcp-artifact", "google_compute_router": "gcp-router",
	"google_compute_router_nat": "gcp-nat", "google_logging_project_sink": "gcp-logging",
}

// categoryKeywords is checked in this fixed priority order; a keyword is a run of
// whole `_`-separated tokens.
var categoryKeywords = []struct {
	category string
	keywords []string
}{
	{"security", []string{"iam", "role", "policy", "kms", "key", "secret", "vault", "security", "firewall", "certificate"}},
	{"database", []string{"db", "database", "sql", "redis", "dynamodb", "cosmos", "cache"}},
	{"network", []string{"vpc", "subnet", "network", "route", "gateway", "lb", "dns", "nat", "peering"}},
	{"storage", []string{"bucket", "storage", "disk", "volume", "efs", "blob"}},
	{"compute", []string{"instance", "virtual_machine", "function", "cluster", "container", "app_service", "lambda"}},
}

var providerFamilies = []struct{ prefix, family string }{
	{"aws_", "aws"}, {"azurerm_", "azure"}, {"google_", "gcp"},
}

// Family returns the ordered icon chain of a resource type: its service family
// (when the type is in the table), its provider family (aws, azure or gcp), its
// category, then generic. The provider argument is reserved for providers whose
// types do not carry a recognizable prefix. The chain is never empty and always ends with generic.
func Family(resourceType, _ string) []string {
	var chain []string
	if svc, ok := serviceFamilies[resourceType]; ok {
		chain = append(chain, svc)
	}
	prefix := resourceType
	if i := strings.Index(resourceType, "_"); i >= 0 {
		prefix = resourceType[:i]
	}
	for _, pf := range providerFamilies {
		if strings.HasPrefix(resourceType, pf.prefix) {
			chain = append(chain, pf.family)
			break
		}
	}
	if cat := category(strings.TrimPrefix(resourceType, prefix+"_")); cat != "" {
		chain = append(chain, cat)
	}
	return append(chain, familyGeneric)
}

// category applies the keyword heuristic to the type name without its provider
// prefix, matching whole tokens only.
func category(name string) string {
	padded := "_" + name + "_"
	for _, c := range categoryKeywords {
		for _, kw := range c.keywords {
			if strings.Contains(padded, "_"+kw+"_") {
				return c.category
			}
		}
	}
	return ""
}
