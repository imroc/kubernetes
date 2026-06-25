all: codegen run

codegen:
	GENERATED_FILE_PREFIX=genschema. ../hack/update-codegen.sh openapi
	mv ../pkg/generated/openapi/genschema.openapi.go ./openapi/genschema.openapi.go
	# Clean up intermediate model_name files generated in staging by openapi-gen
	find ../staging ../cmd ../pkg ../plugin -name 'genschema.model_name.go' -delete

run:
	go run . $(if $(OUTPUT_DIR),-o $(OUTPUT_DIR))

clean:
	rm -rf schemas agones envoygateway gateway-api

link-kubeschemas:
	ln -s ~/dotfiles/kubeschemas schemas

# Clone external API repos at pinned versions (from versions.txt).
# Use `make clone` for initial clone, then `make sync-code` to copy types into staging.
clone: clone-agones clone-envoygateway clone-gateway-api

clone-agones:
	git clone https://github.com/googleforgames/agones.git
	@cd agones && git checkout $$(grep '^agones=' versions.txt | cut -d= -f2)

clone-envoygateway:
	git clone https://github.com/envoyproxy/gateway.git envoygateway
	@cd envoygateway && git checkout $$(grep '^envoygateway=' versions.txt | cut -d= -f2)

clone-gateway-api:
	git clone https://github.com/kubernetes-sigs/gateway-api.git gateway-api
	@cd gateway-api && git checkout $$(grep '^gateway-api=' versions.txt | cut -d= -f2)

# Sync external API type definitions into staging/src.
# After cloning, run this to copy the type files into the k8s staging directory.
sync-code: ag eg gwapi

ag:
	mkdir -p ../staging/src/agones/pkg/apis
	cp -r agones/pkg/apis/scheduling.go ../staging/src/agones/pkg/apis/scheduling.go
	cp -r agones/pkg/apis/allocation ../staging/src/agones/pkg/apis/
	cp -r agones/pkg/apis/agones ../staging/src/agones/pkg/apis/

eg:
	cp ./envoygateway/api/v1alpha1/*_types.go ../staging/src/gateway.envoyproxy.io/v1alpha1/
	cp ./envoygateway/api/v1alpha1/fault_injection.go ../staging/src/gateway.envoyproxy.io/v1alpha1/

gwapi:
	cp -r gateway-api/apis ../staging/src/sigs.k8s.io/gateway-api/
