# Scripting Guide

This guide explains the basics of events available to `brigade.js` files.

For more, see the [Brigade Scripting Guide](https://github.com/Azure/brigade/blob/master/docs/topics/scripting.md)

# Brigade GitLab Events

Brigade listens for certain things to happen, this gateway provides those such events from a GitLab repository. The events that Brigade listens for are configured in your project.

When Brigade observes such an event, it will load the `brigade.js` file and see if there is an event handler that matches the event.

For example:

```javascript
const { events } = require("brigadier")

events.on("push", () => {
  console.log("==> handling an 'push' event")
})
```

The GitLab Gateway produces 8 events:

```
push
tag
issue
comment
mergerequest
wikipage
pipeline
build
```

These are based on the events described in the [GitLab Webhooks API](https://gitlab.com/help/user/project/integrations/webhooks) guide.
Note the *Issue* and *ConfidentialIssue* events are treated the same as an `issue` event.