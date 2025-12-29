resource "google_artifact_registry_repository" "foresee" {
  location      = var.region
  repository_id = "foresee"
  description   = "Docker container images for Foresee services"
  format        = "DOCKER"
}
