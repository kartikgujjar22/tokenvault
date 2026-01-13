const projectMap = {
  5000: "example_project",
};

var currentPort = pm.request.url.port;
var currentHost = pm.request.url.getHost();

var project = projectMap[currentPort];

if (!project) {
  project = projectMap[currentHost];
}

if (!project) {
  console.log(
    "TokenVault: No project mapped for port" +
      currentPort +
      ". Skipping auti-fetch"
  );
} else {
  console.log(
    "TokenVault: Detected Port " +
      currentPort +
      " -> Fetching token for ' " +
      project +
      "'"
  );

  pm.sendRequest(
    {
      url: "http://localhost:9999/fetch/" + project,
      method: "GET",
    },
    function (err, res) {
      if (!err && res.code === 200) {
        var data = res.json();
        pm.environment.set("token", data.token);
        console.log("TokenVault: Token updated for " + project);
      } else {
        console.log("Tokenvault: Token not found for " + project);
      }
    }
  );
}
