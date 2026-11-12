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

function showOrHideProofOfDeathReceivedContent() {
    const currentUrl = window.location.href;
    console.log("in func");
        htmx.addClass(htmx.find("#proof-of-death-receieved-no"), "hide");
        htmx.removeClass(htmx.find("#proof-of-death-receieved-yes"), "hide");
    // } else {
    //     htmx.addClass(htmx.find("#proof-of-death-receieved-no"), "hide");
    //     htmx.removeClass(htmx.find("#proof-of-death-receieved-yes"), "hide");
    // }
}

document.body.addEventListener("htmx:afterOnLoad", function () {
    return showOrHideProofOfDeathReceivedContent();
});

htmx.onLoad((content) => {
    initAll();

    htmx.findAll(content, ".summary").forEach((element) => {
        htmx.on(`#${element.id}`, "click", () =>
            htmx.toggleClass(htmx.find(`#${element.id}-reveal`), "hide"),
        );
    });

    htmx.findAll("#proofofdeath").forEach((el) => {
        console.log("pod");
        console.log(el.valueOf());
        return showOrHideProofOfDeathReceivedContent(el.valueOf());
    });
    htmx.findAll("#proofofdeath-yes").forEach(() => {
        console.log("pod yes");
        return showOrHideProofOfDeathReceivedContent();
    });
    htmx.findAll("#proofofdeath-no").forEach(() => {
        console.log("pod no");
        return showOrHideProofOfDeathReceivedContent();
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