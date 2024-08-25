# Contribution Guide

This section is only for those who want to contribute to the [framework](https://github.com/componego/componego){:target="_blank"}.

<hr/>

Contributions can be new features, changes to existing features, tests, documentation (such as developer guides, examples, or specifications), bug fixes, optimizations, or just good suggestions.

The main rules for making changes to the codebase are as follows:

1. We use [pre-commit](https://pre-commit.com/){:target="_blank"} to run a custom hook when committing code. Please install it.
2. After making changes please run tests using our utility:
       ```text hl_lines="1"
       % make tests
       ```
3. Please do not include third-party packages in the framework codebase.
4. Before creating a pull request or merge to our repository, please do a ^^git rebase^^.
5. The language of communication is English. All commits, comments, tasks, and questions must be written in English.

You can use the following command to quickly create a framework contributor environment:
```shell
curl -sSL https://raw.githubusercontent.com/componego/componego/master/tools/create-contributor-env.sh | sh
```
or
```shell
wget -O - https://raw.githubusercontent.com/componego/componego/master/tools/create-contributor-env.sh | sh
```

<hr/>

If you like our framework, you can tell your friends about it, which will help this project develop further.

[Twitter | X](https://twitter.com/share?url=github.com%2Fcomponego%2Fcomponego){:target="_blank" .md-button .social-share-button }
[LinkedIn](https://www.linkedin.com/sharing/share-offsite/?url=github.com%2Fcomponego%2Fcomponego){:target="_blank" .md-button .social-share-button }
[Facebook](https://www.facebook.com/sharer/sharer.php?u=github.com%2Fcomponego%2Fcomponego){:target="_blank" .md-button .social-share-button }
<hr/>

:octicons-heart-fill-24:{ .heart } Thank you.
