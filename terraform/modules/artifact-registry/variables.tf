variable "location" {
  description = "The location of the repository"
  type        = string
}

variable "repository_id" {
  description = "The repository ID"
  type        = string
}

variable "description" {
  description = "Description of the repository"
  type        = string
  default     = ""
}

variable "format" {
  description = "The format of the repository (DOCKER, MAVEN, NPM, etc.)"
  type        = string
  default     = "DOCKER"
}
