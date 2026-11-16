import { initAll } from "govuk-frontend";
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

    // htmx.findAll(content, "#record-deputy-death-form").forEach((element) => {
    //     console.log("find all elements in rdd form");
    //     htmx.on(`#${element.id}`, "click", () =>
    //         htmx.toggleClass(htmx.find(`proof-of-death-received-reveal`), "hide"),
    //     );
    // });

    htmx.findAll("#proofofdeath-yes").forEach((element) => {
        console.log("event listener pod y");
        element.addEventListener("click", () =>
            htmx.removeClass(htmx.find("#proof-of-death-receieved-yes"), "hide"),
        );
        element.addEventListener("click", () =>
            htmx.removeClass(htmx.find("#dateofdeath-required"), "hide"),
        );
        element.addEventListener("click", () =>
            htmx.addClass(htmx.find("#dateofdeath-optional"), "hide"),
        );
    });
    htmx.findAll("#proofofdeath-no").forEach((element) => {
        console.log("event listener pod n");
        element.addEventListener("click", () =>
            htmx.addClass(htmx.find("#proof-of-death-receieved-yes"), "hide"),
        );
        element.addEventListener("click", () =>
            htmx.addClass(htmx.find("#dateofdeath-required"), "hide"),
        );
        element.addEventListener("click", () =>
            htmx.removeClass(htmx.find("#dateofdeath-optional"), "hide"),
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