import {initAll} from 'govuk-frontend';
import "govuk-frontend/dist/govuk/all.mjs";
import "opg-sirius-header/sirius-header.js";
import htmx from "htmx.org/dist/htmx.esm";

document.body.className += ' js-enabled' + ('noModule' in HTMLScriptElement.prototype ? ' govuk-frontend-supported' : '');
initAll();

window.htmx = htmx
htmx.logAll();
htmx.config.responseHandling = [{code:".*", swap: true}]

const formToggler = (suffix) => {
    return {
        resetAll: resetAll(suffix),
        show: show(suffix),
    }
}

const resetAll = (suffix) => () => {
    htmx.findAll(`[id$="-${suffix}"]`).forEach(element => {
        htmx.addClass(element, "hide");
        const input = element.querySelector("input");
        if (input) {
            input.setAttribute("disabled", "true");
            input.removeAttribute("max");
        }
    });
}

const show = (suffix) => (idName) => {
    document.querySelector(`#${idName}`).removeAttribute("disabled");
    htmx.removeClass(htmx.find(`#${idName}-${suffix}`), "hide")
}

//added temporary solution for the active sub-navigation state whilst the microservice is being built.
const setCurrentSubNavigationLink = (link) => {
    const navigation = link.closest(".moj-sub-navigation");

    if (!navigation) {
        return;
    }

    navigation.querySelectorAll(".moj-sub-navigation__link[aria-current]").forEach((currentLink) => {
        currentLink.removeAttribute("aria-current");
    });

    link.setAttribute("aria-current", "page");
}

const initialiseSubNavigation = () => {
    htmx.findAll(".moj-sub-navigation").forEach((navigation) => {
        const links = Array.from(navigation.querySelectorAll(".moj-sub-navigation__link"));

        if (links.length === 0) {
            return;
        }

        const currentPath = `${window.location.pathname}${window.location.hash}`;
        const matchingLink = links.find((link) => {
            const linkUrl = new URL(link.getAttribute("href"), window.location.origin);

            return `${linkUrl.pathname}${linkUrl.hash}` === currentPath;
        });
        const defaultLink = navigation.querySelector('.moj-sub-navigation__link[aria-current="page"]') || links[0];

        setCurrentSubNavigationLink(matchingLink || defaultLink);

        links.forEach((link) => {
            if (link.dataset.subNavigationInitialised === "true") {
                return;
            }

            link.dataset.subNavigationInitialised = "true";
            link.addEventListener("click", () => setCurrentSubNavigationLink(link));
        });
    });
}

initialiseSubNavigation();

// adding event listeners inside the onLoad function will ensure they are re-added to partial content when loaded back in
htmx.onLoad(() => {
    initAll();
    initialiseSubNavigation();

    htmx.findAll(".moj-banner--success").forEach((element) => {
        element.addEventListener("click", () => {
            htmx.addClass(element, "hide");
            const url = new URL(window.location.href);
            if (url.searchParams.has('success')) {
                url.searchParams.delete('success');
                window.history.replaceState({}, '', url.toString());
            }
        });
    });
});
