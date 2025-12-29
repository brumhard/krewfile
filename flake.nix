{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    {
      overlays.default = final: prev: {
        krewfile = self.packages.${prev.stdenv.hostPlatform.system}.default;
      };

      homeManagerModules = {
        default = self.homeManagerModules.krewfile;
        krewfile = import ./hm-module.nix self;
      };
    }
    // flake-utils.lib.eachDefaultSystem (system: let
      name = "krewfile";
      version = "0.6.3";
      pkgs = import nixpkgs {inherit system;};
    in rec {
      packages = {
        default = packages.${name};
        ${name} = pkgs.buildGoModule {
          pname = name;
          version = version;
          vendorHash = "sha256-Z0H01Ts6RlBFwKgx+9YYAd9kT4BkCBL1mvJsRf2ci5I=";
          src = ./.;

          meta = with pkgs.lib; {
            description = "Helper to declaratively manage krew plugins";
            homepage = "https://github.com/brumhard/krewfile";
            maintainers = with maintainers; [brumhard];
            license = licenses.mit;
          };
        };
      };

      apps.default =
        (flake-utils.lib.mkApp {
          drv = packages.default;
        })
        // {
          meta = packages.default.meta;
        };

      devShells.default = import ./shell.nix {inherit pkgs;};
    });
}
