#!/usr/bin/env python3
"""Build the vendor icons embedded in tf2report from the official icon packs.

The packs are downloaded by hand (they are not committed):
  AWS     https://aws.amazon.com/architecture/icons/
  Azure   https://learn.microsoft.com/azure/architecture/icons/
  Google  https://cloud.google.com/icons

Usage:
  scripts/build-icons.py [--aws DIR] [--azure DIR] [--gcp DIR ...]

Each --gcp may be given several times (core products and category packs).
Output: pkg/report/icons/vendor/<cloud>/<family-id>.svg, one file per family id of
pkg/diagram/family.go. Files are cleaned (comments, XML declaration, CSS classes
inlined as attributes) so that they pass the sanitizer of pkg/report/svg_icons.go.
Standard library only.
"""
import argparse
import pathlib
import re
import sys

OUT = pathlib.Path(__file__).resolve().parent.parent / "pkg" / "report" / "icons" / "vendor"

# family id -> file name (without extension) searched anywhere in the pack.
AWS = {
    "aws-ec2": "Arch_Amazon-EC2_48", "aws-autoscaling": "Arch_Amazon-EC2-Auto-Scaling_48",
    "aws-lambda": "Arch_AWS-Lambda_48", "aws-vpc": "Arch_Amazon-Virtual-Private-Cloud_48",
    "aws-subnet": "Arch_Amazon-Virtual-Private-Cloud_48", "aws-route53": "Arch_Amazon-Route-53_48",
    "aws-cloudfront": "Arch_Amazon-CloudFront_48", "aws-apigw": "Arch_Amazon-API-Gateway_48",
    "aws-elb": "Arch_Elastic-Load-Balancing_48", "aws-igw": "Res_Amazon-VPC_Internet-Gateway_48",
    "aws-nat": "Res_Amazon-VPC_NAT-Gateway_48", "aws-sg": "Res_Amazon-VPC_Network-Access-Control-List_48",
    "aws-s3": "Arch_Amazon-Simple-Storage-Service_48", "aws-efs": "Arch_Amazon-EFS_48",
    "aws-ebs": "Arch_Amazon-Elastic-Block-Store_48", "aws-rds": "Arch_Amazon-RDS_48",
    "aws-dynamodb": "Arch_Amazon-DynamoDB_48", "aws-elasticache": "Arch_Amazon-ElastiCache_48",
    "aws-iam": "Arch_AWS-Identity-and-Access-Management_48", "aws-kms": "Arch_AWS-Key-Management-Service_48",
    "aws-eks": "Arch_Amazon-Elastic-Kubernetes-Service_48", "aws-ecs": "Arch_Amazon-Elastic-Container-Service_48",
    "aws-ecr": "Arch_Amazon-Elastic-Container-Registry_48", "aws-sqs": "Arch_Amazon-Simple-Queue-Service_48",
    "aws-sns": "Arch_Amazon-Simple-Notification-Service_48", "aws-cloudwatch": "Arch_Amazon-CloudWatch_48",
}
# Azure files are named <number>-icon-service-<Name>.svg; the value is <Name>.
AZURE = {
    "azure-vm": "Virtual-Machine", "azure-vmss": "VM-Scale-Sets", "azure-rg": "Resource-Groups",
    "azure-vnet": "Virtual-Networks", "azure-subnet": "Subnet", "azure-nsg": "Network-Security-Groups",
    "azure-nic": "Network-Interfaces", "azure-pip": "Public-IP-Addresses", "azure-lb": "Load-Balancers",
    "azure-storage": "Storage-Accounts", "azure-keyvault": "Key-Vaults", "azure-aks": "Kubernetes-Services",
    "azure-acr": "Container-Registries", "azure-sql": "SQL-Database",
    "azure-postgres": "Azure-Database-PostgreSQL-Server", "azure-cosmos": "Azure-Cosmos-DB",
    "azure-redis": "Cache-Redis", "azure-functions": "Function-Apps", "azure-appservice": "App-Services",
    "azure-dns": "DNS-Zones", "azure-monitor": "Monitor", "azure-identity": "Managed-Identities",
    "azure-servicebus": "Azure-Service-Bus", "azure-appgw": "Application-Gateways",
}
GCP = {
    "gcp-gce": "ComputeEngine-512-color-rgb", "gcp-gke": "GKE-512-color", "gcp-gcs": "Cloud_Storage-512-color",
    "gcp-sql": "CloudSQL-512-color", "gcp-run": "CloudRun-512-color-rgb", "gcp-bigquery": "BigQuery-512-color",
    "gcp-disk": "Hyperdisk-512-color",
}


def index(root, pattern):
    return {p.stem: p for p in pathlib.Path(root).expanduser().rglob(pattern)}


def clean(svg):
    """Inline CSS classes as attributes and drop what the sanitizer would refuse."""
    style = re.search(r"<style[^>]*>(.*?)</style>", svg, re.S)
    if style:
        rules = {}
        for sel, body in re.findall(r"\.([\w-]+)\s*\{([^}]*)\}", style.group(1)):
            rules[sel] = re.findall(r"([\w-]+)\s*:\s*([^;]+);?", body)

        def inline(m):
            attrs = {}
            for c in m.group(1).split():
                for k, v in rules.get(c, []):
                    attrs[k] = v.strip()
            return " ".join(f'{k}="{v}"' for k, v in attrs.items())

        svg = re.sub(r'class="([^"]*)"', inline, svg)
        svg = re.sub(r"<style[^>]*>.*?</style>", "", svg, flags=re.S)
        svg = re.sub(r"<defs>\s*</defs>", "", svg)
    svg = re.sub(r"<\?xml[^>]*\?>", "", svg)
    svg = re.sub(r"<!--.*?-->", "", svg, flags=re.S)
    return re.sub(r">\s+<", "><", svg).strip() + "\n"


def build(cloud, mapping, files, name_of):
    out = OUT / cloud
    out.mkdir(parents=True, exist_ok=True)
    for old in out.glob("*.svg"):
        old.unlink()
    missing = []
    for family, key in sorted(mapping.items()):
        src = files.get(name_of(key))
        if src is None:
            missing.append(family)
            continue
        (out / f"{family}.svg").write_text(clean(src.read_text(encoding="utf-8")), encoding="utf-8")
    print(f"{cloud}: {len(mapping) - len(missing)} icons" + (f", missing: {', '.join(missing)}" if missing else ""))


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--aws")
    ap.add_argument("--azure")
    ap.add_argument("--gcp", action="append", default=[])
    args = ap.parse_args()
    if not (args.aws or args.azure or args.gcp):
        ap.error("give at least one pack")
    if args.aws:
        build("aws", AWS, index(args.aws, "*.svg"), lambda k: k)
    if args.azure:
        files = {re.sub(r"^\d+-icon-service-", "", p.stem): p for p in pathlib.Path(args.azure).expanduser().rglob("*.svg")}
        build("azure", AZURE, files, lambda k: k)
    if args.gcp:
        files = {}
        for d in args.gcp:
            files.update(index(d, "*.svg"))
        build("gcp", GCP, files, lambda k: k)


if __name__ == "__main__":
    sys.exit(main())
