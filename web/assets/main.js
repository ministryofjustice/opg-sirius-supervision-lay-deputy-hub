import {initAll} from 'govuk-frontend';
import "govuk-frontend/dist/govuk/all.mjs";
import "opg-sirius-header/sirius-header.js";
import htmx from "htmx.org/dist/htmx.esm";
import("htmx-ext-response-targets");

document.body.className +=
    " js-enabled" + ("noModule" in HTMLScriptElement.prototype
        ? " govuk-frontend-supported"
        : "");
initAll();

window.htmx = htmx
htmx.logAll();
htmx.config.responseHandling = [{code:".*", swap: true}]

htmx.onLoad((content) => {
    initAll();
    htmx.findAll("#proofofdeath").forEach((element) => {
        console.log("clicked proof of death")
        element.addEventListener(
            "change",
            function () {
                document
                    .querySelector("#proof-of-death-receieved-yes")
                    .setAttribute("hide", "false");
                document
                    .querySelector("#proof-of-death-receieved-no")
                    .setAttribute("hide", "true");
            },
            false,
        );
    });

    htmx.findAll("#proofofdeath-no").forEach((element) => {
        console.log("clicked proof of death no")

        element.addEventListener(
            "change",
            function () {
                document
                    .querySelector("#proof-of-death-receieved-yes")
                    .setAttribute("hide", "true");
                document
                    .querySelector("#proof-of-death-receieved-no")
                    .setAttribute("hide", "false");
            },
            false,
        );
    });
});

function onHomePage() {
    const homePageUrlRegex = new RegExp("^\\/(supervision/deputies\\/)?\\d+\\/*$");
    return homePageUrlRegex.test(location.pathname);
}

function storeBackSessionVars(backIndex, href) {
    if (backIndex !== null && location.href === href) {
        sessionStorage.setItem("backIndex", (parseInt(backIndex) - 1).toString());
    }

    if (backIndex === null || href === null || location.href !== href || onHomePage()) {
        sessionStorage.setItem("backIndex", "-1");
        sessionStorage.setItem("href", location.href);
    }
}

storeBackSessionVars(sessionStorage.getItem("backIndex"), sessionStorage.getItem("href"));