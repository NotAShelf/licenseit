{buildGoModule}:
buildGoModule {
  pname = "licenseit";
  version = "0.1.0";

  src = ../.;

  vendorHash = null;

  ldflags = ["-s" "-w"];
}
