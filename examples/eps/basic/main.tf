resource "hcs_enterprise_project" "eps1" {
  name = "test_eps1"
  project_id = "2e7e928357b14afab35e1ebf7be0f879"
  description = "test_eps1"
}

resource "hcs_enterprise_project" "eps2" {
  name = "test_eps2"
  project_id = "2e7e928357b14afab35e1ebf7be0f879"
  description = "test_eps2"
}

resource "hcs_enterprise_project" "eps3" {
  name = "test_eps3"
  project_id = "2e7e928357b14afab35e1ebf7be0f879"
  description = "test_eps3"
}

data "hcs_enterprise_project" "epslist" {
  depends_on = [
    hcs_enterprise_project.eps1,
    hcs_enterprise_project.eps2
  ]
}

resource "hcs_enterprise_project" "eps4" {
  name = "test_${data.hcs_enterprise_project.epslist.instances.1.id}"
  project_id = "2e7e928357b14afab35e1ebf7be0f879"
  description = "test_${data.hcs_enterprise_project.epslist.instances.1.description}"
}

resource "hcs_enterprise_project" "eps5" {
  name = "test_${data.hcs_enterprise_project.epslist.instances.2.id}"
  project_id = "2e7e928357b14afab35e1ebf7be0f879"
  description = "test_${data.hcs_enterprise_project.epslist.instances.2.description}"
}

resource "hcs_enterprise_project" "eps6" {
  name = "test_${data.hcs_enterprise_project.epslist.instances.3.id}"
  project_id = "2e7e928357b14afab35e1ebf7be0f879"
  description = "test_${data.hcs_enterprise_project.epslist.instances.3.description}"
}

resource "hcs_enterprise_project" "eps7" {
  name = "eps7_${resource.hcs_enterprise_project.eps3.name}"
  project_id = "2e7e928357b14afab35e1ebf7be0f879"
  description = "eps7_${resource.hcs_enterprise_project.eps3.name}"
}

resource "hcs_enterprise_project" "eps8" {
  name = "test111"
  project_id = "b195503ac1a54f7ab08a91ede3f22dcf"
  description = "111"
}