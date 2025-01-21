{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  packages = with pkgs; [
    go_1_22
    earthly
  ];

  shellHook = '''';
}
