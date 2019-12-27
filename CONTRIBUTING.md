# Contributing

First things first, thank you for considering adding your contribution to our provider. The larger our community, the faster the provider can grow to Terraform all the things!

## Before you start

- Install the latest version of the provider and take it for a spin: create, update, and destroy a few resources.
- Identify what you want to change -- is it a bugfix? enahncing an existing resource? adding a new resource?
- Look at the existing issues -- has someone proposed the change you're thinking of?

## Your first contribution

### Steps to contributing

1. Create an issue to track your change, following the guidelines in the issue template
2. Clone the repo
3. Make your change
4. Test your change locally
5. Submit your PR
6. Apply any PR feedback
7. [Admin] Merge PR
8. Wait for the next release
9. Contribution complete! :tada:

## Guidelines

- Code formatting should follow the [Go community standards](https://github.com/golang/go/wiki/CodeReviewComments)
- Favor readability over short, terse code
- Use inline comments to help deciper complex or "unique" code
- Leverage the `logrus` logger to output logging statements - by default, prefer a `Debug` or lower level for most statements
- If in doubt, ask! We welcome incomplete PRs if you're looking for guidance or help!
