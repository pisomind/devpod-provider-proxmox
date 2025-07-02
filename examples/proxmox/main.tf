# ==============================================================================
# Proxmox DevPod Provider Terraform Configuration
# ==============================================================================
# This Terraform configuration creates a Proxmox VM for DevPod development
# environments. It provisions a QEMU virtual machine with cloud-init support,
# configures networking, storage, and SSH access for development purposes.
#
# Prerequisites:
# - Proxmox VE cluster with API access
# - Your VM template of choice
# - Valid API token with VM creation permissions
# - Network configuration (IP range, gateway)
#
# ==============================================================================

terraform {
  required_providers {
    proxmox = {
      source  = "Telmate/proxmox"
      version = "3.0.2-rc01"
    }
  }
}

# ==============================================================================
# Variables
# ==============================================================================

variable "state" {
  description = "Desired state of the VM (running, stopped, etc.)"
  type        = string
  default     = "running"
}

variable "disk_size" {
  description = "Size of the VM disk in GB"
  type        = number
  default     = 100
}

variable "ssh_key" {
  description = "SSH public key for user authentication"
  type        = string
  default     = "invalid"
  sensitive   = true
}

variable "devpod_ssh_key" {
  description = "DevPod SSH public key for automated access"
  type        = string
  default     = "invalid"
  sensitive   = true
}

variable "ci_user" {
  description = "Cloud-init username for the VM"
  type        = string
  default     = "invalid"
}

variable "ci_password" {
  description = "Cloud-init password for the VM user"
  type        = string
  default     = "invalid"
  sensitive   = true
}

variable "ci_ip" {
  description = "Static IP address for the VM (CIDR notation, e.g., 192.168.1.100/24)"
  type        = string
  default     = "invalid"
  
  validation {
    condition     = can(cidrhost(var.ci_ip, 0))
    error_message = "IP address must be in valid CIDR notation (e.g., 192.168.1.100/24)."
  }
}

variable "ci_gateway" {
  description = "Gateway IP address for the VM network"
  type        = string
  default     = "invalid"
}

variable "node_name" {
  description = "Proxmox cluster node name where the VM will be created"
  type        = string
  default     = "devstaginghost1"
}

variable "pm_api_url" {
  description = "Proxmox API URL (e.g., https://proxmox.example.com:8006/api2/json)"
  type        = string
  default     = "invalid"
}

variable "pm_api_token_id" {
  description = "Proxmox API token ID for authentication"
  type        = string
  default     = "invalid"
  sensitive   = true
}

variable "pm_api_token_secret" {
  description = "Proxmox API token secret for authentication"
  type        = string
  default     = "invalid"
  sensitive   = true
}

variable "proxmox_vm_id" {
  description = "Unique VM ID for the Proxmox virtual machine"
  type        = string
  default     = "invalid"
  
  validation {
    condition     = can(tonumber(var.proxmox_vm_id)) && tonumber(var.proxmox_vm_id) > 0
    error_message = "VM ID must be a positive integer."
  }
}

variable "proxmox_template_name" {
  description = "Name of the Proxmox VM template to clone from"
  type        = string
  default     = "ubuntu-noble-devbox-base"
}

# ==============================================================================
# Provider Configuration
# ==============================================================================

# Configure the Proxmox provider for VM management
# Requires API token authentication and handles TLS settings
provider "proxmox" {
  pm_api_url          = var.pm_api_url
  pm_api_token_id     = var.pm_api_token_id
  pm_api_token_secret = var.pm_api_token_secret
  pm_tls_insecure     = true  # Set to false for production with valid certificates
  pm_debug            = false # Set to true for debugging API calls
}

# ==============================================================================
# Resources
# ==============================================================================

# Create a Proxmox QEMU virtual machine for DevPod development
# This resource provisions a full-featured development environment with:
# - 4 CPU cores and 16GB RAM
# - Configurable disk storage
# - Cloud-init for automated setup
# - Network configuration with static IP
# - SSH key authentication
resource "proxmox_vm_qemu" "devpod" {
  # Basic VM identification
  vmid = var.proxmox_vm_id
  name = "${var.ci_user}-devbox"
  desc = "DevPod development environment for ${var.ci_user}"

  # Proxmox cluster configuration
  # Node name must match the cluster node name (may not include FQDN)
  target_node = var.node_name

  # Template configuration
  # Clone from the specified Ubuntu Noble DevBox base template
  clone = var.proxmox_template_name

  # VM agent and OS configuration
  agent   = 1            # Enable QEMU agent for enhanced management
  os_type = "cloud-init" # Use cloud-init for automated configuration

  # CPU configuration
  # Optimized for development workloads
  cpu {
    cores   = 4      # Number of CPU cores
    sockets = 1      # Number of CPU sockets
    type    = "host"
  }

  # Memory configuration (16GB)
  memory = 16384
  
  # SCSI controller type for storage
  scsihw = "virtio-scsi-pci"

  # Storage configuration
  disks {
    # IDE controller for cloud-init configuration
    ide {
      ide2 {
        cloudinit {
          storage = "local-lvm"
        }
      }
    }
    
    # SCSI controller for main storage
    scsi {
      scsi0 {
        disk {
          size      = var.disk_size # Configurable disk size
          cache     = "writeback"   # Write cache for better performance
          storage   = "local-lvm"   # Storage backend
          replicate = true          # Enable replication if configured
        }
      }
    }
  }

  # VGA configuration for console access
  vga {
    type   = "std" # Standard VGA adapter
    memory = 4     # Video memory in MB
  }

  # Network interface configuration
  network {
    id     = 0        # Network interface ID
    model  = "virtio" # Virtio network driver for performance
    bridge = "vmbr0"  # Bridge interface on Proxmox host
    # tag = 256       # Uncomment to use VLAN tagging
  }

  # Serial console configuration
  serial {
    id   = 0        # Serial port ID
    type = "socket" # Socket type for console access
  }

  # Boot configuration
  boot = "order=scsi0" # Boot from SCSI disk

  # Network configuration via cloud-init
  # IP address must be in CIDR notation (e.g., 192.168.1.100/24)
  ipconfig0 = "ip=${var.ci_ip},gw=${var.ci_gateway}"

  # SSH key configuration
  # Combines user SSH key and DevPod SSH key for access
  sshkeys = <<EOF
    ${var.ssh_key}
    ${var.devpod_ssh_key}
    EOF

  # Cloud-init user configuration
  ciuser     = var.ci_user     # Username for the VM
  cipassword = var.ci_password # Password for the VM user
}

# ==============================================================================
# Outputs
# ==============================================================================

output "public_ip" {
  description = "The public IP address of the created DevPod VM"
  value       = var.ci_ip
}