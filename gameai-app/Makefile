mock:
	docker run -i -p 8000:8000 -v $(pwd)/../backend/api/swagger.yml:/data/swagger.yaml -t palo/swagger-api-mock

src/api-types.ts: ../backend/api/swagger.yml
	npx sw2dts $< > $@

.PHONY: mock