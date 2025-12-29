{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  packages = with pkgs; [
    go
    earthly
    goreleaser
  ];

  shellHook = '''';
}
