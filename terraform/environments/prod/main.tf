data "google_project" "project" {
  project_id = var.project_id
}

# Artifact Registry repository for Docker images
module "artifact_registry" {
  source = "../../modules/artifact-registry"

  location      = var.region
  repository_id = "foresee"
  description   = "Docker container images for Foresee services"
  format        = "DOCKER"
}
