ENV=dev
COMMAND=""
#############################
# command
#############################
.PHONY: terraform.command
terraform.command:
	@TF_DATA_DIR="env/${ENV}/.terraform" terraform ${COMMAND}
#############################
# init
#############################
.PHONY: terraform.init.dev
terraform.init.dev:
	@make terraform.command ENV=dev COMMAND="init -upgrade -backend-config env/dev/backend.tfvars"
.PHONY: terraform.init.prd
terraform.init.prd:
	@make terraform.command ENV=prd COMMAND="init -upgrade -backend-config env/prd/backend.tfvars"
#############################
# plan
#############################
.PHONY: terraform.plan.dev
terraform.plan.dev:
	@make terraform.command ENV=dev COMMAND="plan -var-file env/dev/terraform.tfvars"
.PHONY: terraform.plan.prd
terraform.plan.prd:
	@make terraform.command ENV=prd COMMAND="plan -var-file env/prd/terraform.tfvars"
#############################
# apply
#############################
.PHONY: terraform.apply.dev
terraform.apply.dev:
	@make terraform.command ENV=dev COMMAND="apply -var-file env/dev/terraform.tfvars"
.PHONY: terraform.apply.prd
terraform.apply.prd:
	@make terraform.command ENV=prd COMMAND="apply -var-file env/prd/terraform.tfvars"