import flask
import os
import pulumi
from pulumi import automation as auto
from jaxxstorm_pulumi_productionapp import Deployment, DeploymentArgs

# references the templates and static files in the assets dir
template_dir = os.path.abspath('./assets/templates')
css_dir = os.path.abspath('./assets/static')
app = flask.Flask(__name__, template_folder=template_dir, static_folder=css_dir)


app.secret_key = "super-secret-key"


# we want all our deployments to go into the same stack
project_name = "productionapp-platyform"


# we use the component resource we built earlier as our Pulumi program
def create_pulumi_program(name: str, image: str, port: int):
    app = Deployment(
        name, DeploymentArgs(image=image, port=port)
    )


@app.route("/ping", methods=["GET"])
def ping():
	return flask.jsonify("pong!", 200)


@app.route("/", methods=["GET"])
def list_deployments():
	deployments = []
	try:
		ws = auto.LocalWorkspace(project_settings=auto.ProjectSettings(name=project_name, runtime="python"))
		all_stacks = ws.list_stacks()
		for stack in all_stacks:
			stack = auto.select_stack(stack_name=stack.name,
			                          project_name=project_name,
			                          # no-op program, just to get outputs
			                          program=lambda: None)
			outs = stack.outputs()
			deployments.append({"name": stack.name})
	except Exception as exn:
		flask.flash(str(exn), category="danger")
	return flask.render_template("index.html", deployments=deployments)

@app.route("/new", methods=["GET", "POST"])
def create_deployment():
    """creates new deployment"""
    if flask.request.method == "POST":
        stack_name = flask.request.form.get("name")
        image = flask.request.form.get("image")
        port = flask.request.form.get("port")

        def pulumi_program():
            return create_pulumi_program(stack_name, image, int(port))

        try:
            # create a new stack, generating our pulumi program on the fly from the POST body
            stack = auto.create_stack(
                stack_name=str(stack_name),
                project_name=project_name,
                program=pulumi_program,
            )
            # deploy the stack, tailing the logs to stdout
            stack.up(on_output=print)
            flask.flash(f"Successfully created deployment '{stack_name}'", category="success")
        except auto.StackAlreadyExistsError:
            flask.flash(
                f"Error: Deployment with name '{stack_name}' already exists, pick a unique name",
                category="danger",
            )

        return flask.redirect(flask.url_for("list_deployments"))

    return flask.render_template("create.html")

@app.route("/<string:id>/delete", methods=["POST"])
def delete_deployment(id: str):
    stack_name = id
    try:
        stack = auto.select_stack(stack_name=stack_name,
                                  project_name=project_name,
                                  # noop program for destroy
                                  program=lambda: None)
        stack.destroy(on_output=print)
        stack.workspace.remove_stack(stack_name)
        flask.flash(f"Deployment '{stack_name}' successfully deleted!", category="success")
    except auto.ConcurrentUpdateError:
        flask.flash(f"Error: Deployment '{stack_name}' already has update in progress", category="danger")
    except Exception as exn:
        flask.flash(str(exn), category="danger")

    return flask.redirect(flask.url_for("list_deployments"))

if __name__ == "__main__":
    app.run(host="localhost", port=5050, debug=True)