# Fork
Fork fist to your repository. `https://github.com/<your-username>/fist`
# Clone
Clone `https://github.com/<your-username>/fist` your own repository to develop locally.
# Set remote upstream
```
git remote add upstream https://github.com/sealyun/test.git
git remote set-url --push upstream no-pushing
```
# Fetch from upstream
Merge fanux/fist to `<your-username>/fist`
```
git pull upstream master
git push
```

Merge `<your-username>/fist` to fanux/fist just using github Pull requests.

# Workflow
```
                         Pull request
   upstream(fanux/fist)<-------------------------origin(<you>/fist)
         |                                        ^
         | pull upstream master            push   |
         +---------------------->localGit---------+
```
