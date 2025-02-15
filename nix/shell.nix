{
  mkShell,
  go,
  gopls,
  delve,
}:
mkShell {
  name = "licenseit";
  packages = [
    delve
    go
    gopls
  ];
}
