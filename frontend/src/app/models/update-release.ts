export type ReleaseEventPayload = GithubRelease;

interface GithubRelease {
  tag_name: string;
  body: string;
	assets: {
    name: string;
    browser_download_url: string;
  }[];
	created_at: string;
}

export class UpdateRelease {
  Version: string;
  Changelog: string;
	Assets: {
    Name: string;
    Url: string;
  }[];
	CreatedAt: string;
  OldVersion: string;

  constructor(json: GithubRelease) {
    this.Version = json.tag_name;
    this.Changelog = json.body;
    this.Assets = json.assets.map(asset => ({ Name: asset.name, Url: asset.browser_download_url }));
    this.CreatedAt = json.created_at;
  }

  static fromJSON(json: GithubRelease): UpdateRelease {
    return new this(json)
  }
}
