{ lib, buildGoModule }:

# TODO: somehow include krew as dependency?
buildGoModule rec {
    pname = "krewfile";
    version = "dev";
    vendorSha256 = "sha256-Z0H01Ts6RlBFwKgx+9YYAd9kT4BkCBL1mvJsRf2ci5I=";
    src = ./.;

    # tests expect the source files to be a build repo
    doCheck = false;

    meta = with lib; {
        description = "helper to declaratively manage krew plugins";
        homepage = "https://goreleaser.com";
        maintainers = with maintainers; [ brumhard ];
        license = licenses.mit;
    };
}
