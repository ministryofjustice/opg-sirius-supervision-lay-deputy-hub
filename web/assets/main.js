import {initAll} from 'govuk-frontend';
import "govuk-frontend/dist/govuk/all.mjs";
import "opg-sirius-header/sirius-header.js";
import htmx from "htmx.org/dist/htmx.esm";

document.body.className += ' js-enabled' + ('noModule' in HTMLScriptElement.prototype ? ' govuk-frontend-supported' : '');
initAll();

window.htmx = htmx
htmx.logAll();
htmx.config.responseHandling = [{code:".*", swap: true}]

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