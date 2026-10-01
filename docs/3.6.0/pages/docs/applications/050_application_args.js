<doc-view>

<h2 id="_set_application_arguments">Set Application Arguments</h2>
<div class="section">
<p>When running a custom application there may be a requirement to pass arguments to the application&#8217;s main class.
By default, there are no application arguments but any arguments required can be specified in the <code>application.args</code> list
in the <code>Coherence</code> resource spec.</p>

<p>The <code>application.args</code> is a list of string values, each value in the list is passed as an argument, in the order
that they are specified in the list.</p>

<p>For example, a deployment uses a custom image <code>catalogue:1.0.0</code> that requires a custom main class
called <code>com.acme.Catalogue</code>, and that class takes additional arguments.
In this example we&#8217;ll use two fictitious arguments such as a name and a language for the catalogue.
the <code>Coherence</code> resource would look like this:</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: catalogue:1.0.0
  application:
    main: com.acme.Catalogue <span class="conum" data-value="1" />
    args:                    <span class="conum" data-value="2" />
      - "--name=Books"
      - "--language=en_GB"</markup>

<ul class="colist">
<li data-value="1">The <code>com.acme.Catalogue</code> will be run as the main class.</li>
<li data-value="2">The arguments passed to the <code>com.acme.Catalogue</code> class will be <code>--name=Books</code> and <code>--language=en_GB</code></li>
</ul>
<p>The example would be equivalent to the Coherence container running:</p>

<markup
lang="bash"

>$ java com.acme.Catalogue --name=Books --language=en_GB</markup>

</div>

<h2 id="_environment_variable_expansion">Environment Variable Expansion</h2>
<div class="section">
<p>The Operator supports environment variable expansion in program arguments.
The runner in the Coherence container will replace <code>${var}</code> or <code>$var</code> in the program arguments with the corresponding environment variable name.</p>

<markup
lang="yaml"

>apiVersion: coherence.oracle.com/v1
kind: Coherence
metadata:
  name: test
spec:
  image: catalogue:1.0.0
  application:
    main: com.acme.Catalogue
    args:
      - "${HOSTNAME}"  <span class="conum" data-value="1" /></markup>

<ul class="colist">
<li data-value="1">The argument passed to the <code>com.acme.Catalogue</code> main method will resolve to the value of the <code>HOSTNAME</code> environment variable.</li>
</ul>
<p>Any environment variable that is present when the Coherence container starts can be used, this would include variables created as part of the image and variables specified in the Coherence yaml.</p>

</div>
</doc-view>
