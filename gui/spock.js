new Vue({
    el: "#spock",
    data: {
        endpoint: "",
        checks: []
    },
    created: function() {
        this.fetchChecks()
    },
    methods: {
        bg: function(check, p) {
            v = _.at(check, p)[0]
            if (v == true || v == "checked:ok") {
                return "ok-bg"
            } else {
                return "fail-bg"
            }
        },
        ago: function(check, p) {
            d = _.at(check, p)[0]
            return moment(d).fromNow()
        },
        property: function(check, p) {
            return _.at(check, p)[0]
        },
        fetchChecks: function() {
            app = this
            app.checks = []
            fetch(app.endpoint + "/_all")
            .then(function (response) { return response.json() })
            .then(function (checks) { app.checks = _.values(checks) })
            .catch(function () { 
                app.endpoint = prompt("Spock Endpoint", "http://localhost:8080")
                app.fetchChecks()
            })
        }
    }
})
