{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        version = "0.1.1";
        name = "krewfile";
      in
      rec {
        packages = {
          default = pkgs.buildGoModule {
            pname = name;
            version = version;
            vendorSha256 = "sha256-Z0H01Ts6RlBFwKgx+9YYAd9kT4BkCBL1mvJsRf2ci5I=";
            src = ./.;

            meta = with pkgs.lib; {
              description = "Helper to declaratively manage krew plugins";
              homepage = "https://goreleaser.com";
              maintainers = with maintainers; [ brumhard ];
              license = licenses.mit;
            };
          };
          ociImage = pkgs.dockerTools.buildLayeredImage {
            name = "ghcr.io/brumhard/${name}";
            tag = version;
            contents = [ packages.default ];
            config = {
              Entrypoint = [ "${packages.default}/bin/${name}" ];
              Cmd = [ "-help" ];
            };
          };
        };

        apps = {
          default = flake-utils.lib.mkApp {
            drv = packages.default;
          };
        };

        devShell = import ./shell.nix { inherit pkgs; };
      });
}
