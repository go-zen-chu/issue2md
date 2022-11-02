---
title: test issue
url: https://github.com/go-zen-chu/issue2md/issues/2
labels: [invalid,test]
---
This is a test issue for checking issue2md works as expected
github action should run
yes it runs

![image](https://user-images.githubusercontent.com/1454332/196034659-04a85305-b238-4490-8523-f6c7fa69c60d.png)

next try

→ started :)

![image](https://user-images.githubusercontent.com/1454332/196034766-ccb3cf1b-385d-47fb-8bcd-bd6209317144.png)

it should work
worked as expected :) jolly good
![image](https://user-images.githubusercontent.com/1454332/197346843-652a1a76-cf76-4ad1-9013-fb2df57361b4.png)

testing if you can get issue number and title
result. Seems it worked for issue title.

https://github.com/go-zen-chu/issue2md/actions/runs/3305906536/jobs/5456396416

![image](https://user-images.githubusercontent.com/1454332/197375337-47f30d94-2f80-4dbf-b2a8-d735d5bab495.png)

try again with labels
by using join function you can join labels array with comma
https://docs.github.com/en/actions/learn-github-actions/expressions#example-4

```
          # join function from https://docs.github.com/en/actions/learn-github-actions/expressions#example-4
          issue-labels: ${{ join(github.event.issue.labels.*.name,',') }}
```
finally we got labels with csv string!
![image](https://user-images.githubusercontent.com/1454332/197384129-f9dafcc5-524e-4a8e-899c-203604fc3e63.png)

- test when no change detected → worked as expected

![image](https://user-images.githubusercontent.com/1454332/199508391-5e92fdde-92c9-4b89-9a78-f69406a05722.png)

- test when there is change in issue
