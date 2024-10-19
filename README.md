# Gommit

Gommit is a powerful tool for enforcing consistent commit message formats in your Git repositories. It helps teams maintain clean and informative commit histories by validating commit messages against predefined rules.

## Features

- Enforces the [Conventional Commits](https://www.conventionalcommits.org/) specification
- Automatically sets up in repositories with zero configuration for developers
- Cross-platform support (Linux, macOS, Windows)
- Easy integration for repository maintainers

## For Repository Maintainers

### Integrating Gommit into Your Repository

1. Download the latest `gommit-integration-<os>-<arch>` binary for your system from the [releases page](https://github.com/Moukrea/gommit/releases).

2. Place the binary in the root of your repository and run it:

   ```
   ./gommit-integration-<os>-<arch>
   ```

3. The integration tool will set up all necessary files and scripts in your repository.

4. Commit and push the changes:

   ```
   git add .
   git commit -m "chore: integrate Gommit for commit message validation"
   git push
   ```

5. Gommit is now integrated into your repository. All developers who clone or pull from this repository will automatically have Gommit set up for them.

For those who prefer to set up Gommit manually or want to understand the integration process in detail, please refer to our [Manual Integration Guide](docs/integration.md).

### What Gets Added to Your Repository

The integration tool adds the following to your repository:

- `.gommit/`: A directory containing Gommit downloader binaries for various platforms
- `scripts/`: A directory containing setup scripts and Git hooks
- Git hooks: `commit-msg`, `post-checkout`, and `post-merge`

These additions ensure that Gommit works seamlessly for all developers, regardless of their operating system.

## For Developers

If you're working on a repository that has Gommit integrated, you don't need to do anything special. Gommit will be automatically set up when you clone the repository or switch branches.

### Writing Commit Messages

When you make a commit, Gommit will validate your commit message. A valid commit message should follow this format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Where `<type>` is one of:

- feat: A new feature
- fix: A bug fix
- docs: Documentation only changes
- style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
- refactor: A code change that neither fixes a bug nor adds a feature
- perf: A code change that improves performance
- test: Adding missing tests or correcting existing tests
- chore: Changes to the build process or auxiliary tools and libraries such as documentation generation

For example:

```
feat(user-auth): add password reset functionality

- Implement password reset form
- Add email sending service for reset links
- Update user model with reset token field

Closes #123
```

If your commit message doesn't meet the required format, Gommit will prevent the commit and provide feedback on what needs to be corrected.

## Troubleshooting

If you encounter any issues with Gommit:

1. Ensure you have the latest version of Gommit integrated into your repository.
2. Check that the `.gommit` directory and `scripts` directory are present in your repository root.
3. Verify that Git hooks are properly set up by running `git config core.hooksPath`. It should point to the `scripts/git-hooks` directory in your repository.
4. If Gommit isn't running on commit, try running `./scripts/setup-hooks.sh` from your repository root.

## Contributing

Contributions to Gommit are welcome! Please see our [Contributing Guidelines](CONTRIBUTING.md) for more details.

## License

Gommit is released under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

If you encounter any problems or have any questions, please [open an issue](https://github.com/Moukrea/gommit/issues/new) on our GitHub repository.

Thank you for using Gommit to keep your commit history clean and informative!

---
*Note: This project was developed with assistance from AI tools (aider.chat and Claude 3.5 Sonnet).*