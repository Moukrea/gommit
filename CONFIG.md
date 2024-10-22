# Gommit Configuration Guide

Gommit allows you to customize its behavior using a configuration file. This guide explains how to create, use, and customize the configuration file to set rules for your commit messages.

## Creating a Configuration File

To create a Gommit configuration file:

1. Create a file named `gommit.conf.yaml` in the same directory as the Gommit executable.
2. Alternatively, you can create a `.gommit/gommit.conf.yaml` file in your project's root directory.
3. Open the file in your preferred text editor.

## Basic Configuration Structure

The configuration file uses YAML format. Here's the basic structure:

```yaml
disabled_rules:
  - rule_name_1
  - rule_name_2
header_max_length: 50
body_line_max_length: 72
allowed_types:
  - type1
  - type2
  - type3
```

## Available Rules

Gommit comes with several built-in rules:

- `header-format`: Header must be in format: <type>[optional scope][!]: <description>
- `header-max-length`: Header must not exceed the configured max length
- `header-lowercase`: Header (short description) must be all lowercase
- `description-case`: Description must start with lowercase
- `body-line-max-length`: Body lines must not exceed the configured max length
- `footer-format`: Footer must be in format: <token>: <value>
- `breaking-change`: Breaking changes must be indicated in footer
- `auto-breaking-change`: Automatically add BREAKING CHANGE to footer when '!' is present in header
- `type-enum`: Type must be one of the allowed types
- `type-case`: Type must be in lowercase
- `type-empty`: Type must not be empty
- `scope-case`: Scope must be in lowercase
- `subject-empty`: Subject must not be empty

## Customizing Rules

To customize the configuration, you can:

1. Disable specific rules by adding them to the `disabled_rules` list.
2. Set the `header_max_length` and `body_line_max_length`.
3. Define the `allowed_types` for commit messages.

For example:

```yaml
disabled_rules:
  - header-lowercase
  - scope-case
header_max_length: 60
body_line_max_length: 80
allowed_types:
  - feat
  - fix
  - docs
  - style
  - refactor
  - perf
  - test
  - build
  - ci
  - chore
  - revert
```

This configuration:
- Disables the `header-lowercase` and `scope-case` rules
- Sets the maximum header length to 60 characters
- Sets the maximum body line length to 80 characters
- Defines the allowed commit types

## Default Configuration

If no configuration file is found, Gommit uses the following default settings:

```yaml
header_max_length: 50
body_line_max_length: 72
allowed_types:
  - feat
  - fix
  - docs
  - style
  - refactor
  - perf
  - test
  - build
  - ci
  - chore
  - revert
```

## Using the Configuration File

Once you've created and customized your configuration file:

1. Place it in the appropriate location (next to the Gommit executable or in your project's `.gommit` directory).
2. Gommit will automatically use this configuration for all future commits.

## Updating the Configuration

To update the configuration:

1. Edit the configuration file.
2. Save the changes.

The new configuration will apply to all subsequent commits.

Remember, the goal of customizing Gommit is to find the right balance between enforcing good practices and accommodating your team's workflow. Regularly review and adjust your configuration as needed.
