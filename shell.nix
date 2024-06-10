{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  packages = with pkgs; [
    go

    # for terminal gifs
    nodejs
    nodePackages.npm
    pv
    asciinema
    (pkgs.writeShellScriptBin "gen-term-animations" ''
      npx svg-term \
        --command="bash ./ci/type-command.sh" \
        --out docs/term-animation.svg \
        --window=true \
        --width=80
    '')
  ];

  shellHook = '''';
}
