# ghfollow

This bot looks at your timeline and try to find Follow event. If one of your friend follow to somebody, bot will to

### Disclaimer

You run this application at own risk

The more subscribers there are, the longer it takes to find new ones

### Running

You need this environment variables: 

```sh
GITHUB_TOKEN= # GitHub token (PAT) from https://github.com/settings/tokens

GITHUB_RSS= # RSS feed from your home Timeline (url of Subscribe to your news feed in the bottom of page), for example: https://github.com/esin.private.atom?token=ABCDE 

GITHUB_USERNAME= # Your username
```

Or you can run Docker image like:
```sh
docker run -ti -e GITHUB_TOKEN=ghp_xxxx -e GITHUB_RSS='https://github.com/esin.private.atom?token=ABCDE' -e GITHUB_USERNAME='YourUsername' es1n/ghfollow
```

This project should be private, but something goes wrong
