self:
{ config, lib, pkgs, ... }:
with lib;
let
  cfg = config.programs.krewfile;
  finalPackage = self.packages.${pkgs.system}.krewfile;
in
{
  options.programs.krewfile = {
    enable = mkEnableOption "krewfile";
    installKrew = mkEnableOption "installKrew";
  };
  config = mkIf cfg.enable {
    home.packages = ([
      finalPackage
    ] ++ optionals cfg.installKrew [ pkgs.krew ]);
  };
}
