{
  lib,
  buildGoModule,
  tmux,
  ghq,
}:
buildGoModule {
  pname = "tmux-githop";
  version = "0.1.0";

  src = ../.;

  vendorHash = "sha256-xTZBpNc7P8jLCPpsv3cTclVNtLwbW3O/LcF4mYAzXsM=";

  nativeBuildInputs = [tmux ghq];

  subPackages = ["cmd/tmux-githop"];

  checkPhase = ''
    go test ./...
  '';

  meta = {
    description = "Fast tmux session hopping between git repos";
    homepage = "https://github.com/callumio/tmux-githop";
    license = lib.licenses.mit;
    maintainers = with lib.maintainers; [callumio];
  };
}
