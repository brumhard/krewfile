self:
{
  config,
  lib,
  pkgs,
  ...
}:

with lib;

let
  cfg = config.programs.krewfile;
  finalPackage = self.packages.${pkgs.system}.krewfile;
  krewfileContent = pkgs.writeText "krewfile" (
    (concatStringsSep "\n" (map (key: "index ${key} ${getAttr key cfg.indexes}") (attrNames cfg.indexes))) + "\n\n" + (concatStringsSep "\n" cfg.plugins)
  );

  args = if cfg.upgrade then "-upgrade" else "";
in
{
  options.programs.krewfile = {
    enable = mkEnableOption "krewfile";

    krewPackage = mkOption {
      type = types.package;
      default = pkgs.krew;
      defaultText = literalExpression "pkgs.krew";
      description = "Krew package to install.";
    };

    krewRoot = mkOption {
      type = types.path;
      default = "${config.home.homeDirectory}/.krew";
      description = "Path where all krew-related files will be installed and stored.";
    };

    plugins = mkOption {
      type = with types; listOf str;
      default = [ ];
      defaultText = literalExpression "[ "edit-status" ]";
      description = "List of plugins to be installed.";
    };

    indexes = mkOption {
      type = with types; attrsOf str;
      default = { };
      defaultText = ''{ netshoot = "https://github.com/nilic/kubectl-netshoot.git" }'';
      description = "List of extra indexes to be added, where key is index name, and value is index URL";
    };

    upgrade = mkOption {
      type = types.bool;
      default = false;
      defaultText = literalExpression "false";
      description = "Enable auto update of plugins.";
    };
  };

  config = mkIf cfg.enable {
    home.packages = [ finalPackage cfg.krewPackage ];
    home.extraActivationPath = [ pkgs.git ];

    home.sessionVariables.KREW_ROOT = "${cfg.krewRoot}";
    home.sessionPath = [ "${cfg.krewRoot}/bin" ];

    home.activation.krew = hm.dag.entryAfter [ "installPackages" ] ''
      export KREW_ROOT="${cfg.krewRoot}";

      run ${finalPackage}/bin/${finalPackage.pname} \
        -command ${cfg.krewPackage}/bin/${cfg.krewPackage.pname} \
        -file ${krewfileContent} ${args}
    '';
  };
}
