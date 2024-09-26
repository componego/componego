# Contribution Guide

This section is intended for users who wish to contribute to the [framework](https://github.com/componego/componego){:target="_blank"}.

<hr/>

Contributions can include new features, changes to existing features, tests, documentation (such as developer guides, examples, or specifications), bug fixes, optimizations, or valuable suggestions.

Please follow the main rules for making changes to the codebase:

1. We use [pre-commit](https://pre-commit.com/){:target="_blank"} to run a custom hook when code is committed. Please install it.
2. After making changes, please run the tests using our utility:
       ```text hl_lines="1"
       % make tests
       ```
3. Please do not include third-party packages in the framework's codebase.
4. Before creating a pull request or merging to our repository, please perform a ^^git rebase^^.
5. The language of communication is English. All commits, comments, tasks, and questions must be written in English.

You can use the following commands to quickly create a framework contributor environment:
```shell
curl -sSL https://raw.githubusercontent.com/componego/componego/master/tools/create-contributor-env.sh | sh
```
or
```shell
wget -O - https://raw.githubusercontent.com/componego/componego/master/tools/create-contributor-env.sh | sh
```

<hr/>

If you like our framework, you can share it with your friends, which will help the project develop further.

[Twitter | X](https://twitter.com/share?url=github.com%2Fcomponego%2Fcomponego){:target="_blank" .md-button .social-share-button }
[LinkedIn](https://www.linkedin.com/sharing/share-offsite/?url=github.com%2Fcomponego%2Fcomponego){:target="_blank" .md-button .social-share-button }
[Facebook](https://www.facebook.com/sharer/sharer.php?u=github.com%2Fcomponego%2Fcomponego){:target="_blank" .md-button .social-share-button }
<hr/>

:octicons-heart-fill-24:{ .heart } Thank you.
