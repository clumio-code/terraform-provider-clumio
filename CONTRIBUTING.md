## Contributing

All contributions are welcome. See the [README](README.md) to build the provider and set up the
acceptance test environment.

### Verified Commits

Make sure that you sign all your commits. See [Github - Signing Commits](https://help.github.com/articles/signing-commits/) for more details.

### Repository Layout

The provider uses the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
and wraps [clumio-go-sdk](https://github.com/clumio-code/clumio-go-sdk).

- `clumio/plugin_framework/provider.go`: registers all resources and data sources.
- `clumio/plugin_framework/clumio_<name>/`: one package for each resource or data source.
- `clumio/sdk_clients/`: thin interfaces over the SDK, with mockery mocks.
- `clumio/plugin_framework/common/`: shared helpers, for example `ApiClient`, `PollTask` and API error parsing.
- `examples/`, `templates/`, `docs/`: examples and the generated registry docs.

A resource package splits the Terraform plumbing from the API logic:

| File | Content |
|------|---------|
| `resource_<name>.go` | Framework methods: `Metadata`, `Configure`, CRUD, `ImportState` |
| `resource.go` | SDK calls and model population |
| `schema.go` | Model struct (`tfsdk` tags) and schema |
| `const.go` | Schema attribute names |
| `resource_test.go` | Unit tests (build tag `unit`) |
| `resource_<name>_test.go` | Acceptance tests (build tag `basic`, `post_process`, `sso`, `gcp`, ...) |

Data sources use the same split with `data_source_<name>.go`, `data_source.go` and
`data_source_schema.go`.

### Add a Resource

1. If `clumio/sdk_clients/` has no client for the API, add the interface in
   `clumio/sdk_clients/<name>.go`. Then add it to `.mockery.yaml` (keep the list sorted) and run
   `make mockery`.
2. Add the package `clumio/plugin_framework/clumio_<name>/` with the files above.
3. Register the resource in `Resources()` (or `DataSources()`) in `provider.go`.
4. Add `examples/resources/clumio_<name>/resource.tf` and `import.sh`.
5. Run `go generate` to update the docs.
6. Add unit tests and acceptance tests.

### Conventions

Schema and state:

- Keep `tfsdk` tags equal to the `const.go` names.
- Use `RequiresReplace()` on fields that the API cannot update.
- Use `UseStateForUnknown()` on computed fields that do not change.
- Wait for async operations with `common.PollTask`.

Errors:

- Report failures as Terraform diagnostics, close to the SDK call.
- Use `common.ParseMessageFromApiError` for the detail. It also normalizes 401 and 403.
- Put the operation, the resource name and the ID (if known) in the summary.
- For a `nil` SDK response, use `common.NilErrorMessageSummary` and `common.NilErrorMessageDetail`.
- On 404 in `Read`, remove the resource from state. On 404 in `Delete`, treat it as deleted.
- After each `Append`, return if `HasError()` is true.

Security and logging:

- Do not log API tokens, credentials, sensitive attributes, or full API bodies.
- Use `tflog` only to debug a Terraform run. Diagnostics must explain the failure alone.

Context:

- Pass the framework `context.Context` through. Do not use `context.TODO()` or
  `context.Background()` in provider code.

### Testing

Unit tests mock the SDK clients and must pass for every change. Test these cases:

- An SDK API error and a `nil` SDK response.
- Not-found in `Read` and `Delete`.
- Optional attributes that stay null or empty.
- Computed values from the SDK response.
- Task polling success and failure, if the resource uses tasks.

Use the mockery `EXPECT()` API with `Times(1)`. Use `mock.Anything` only when the request is not
relevant to the case.

A new resource must have acceptance tests. A change to a resource must keep its acceptance tests
passing.

Run one package:

```shell
gotestsum --format=pkgname-and-test-fails -- -vet=off ./clumio/plugin_framework/clumio_policy/ -v -tags="unit" -run TestCreatePolicy -timeout 90s
TF_ACC=1 gotestsum --format=pkgname-and-test-fails -- -vet=off ./clumio/plugin_framework/clumio_policy/ -v -tags="basic" -run TestAccClumioPolicy -timeout 120m
```

### Generate Docs

Run `go generate` after a change to a schema, an example or a template. The directory name must be
`terraform-provider-clumio`, or `tfplugindocs` cannot find the provider.

### Release Version

Change the version in two places together:

- `VERSION` in `GNUmakefile`
- `clumioTfProviderVersionHeaderValue` in `clumio/plugin_framework/const.go`
