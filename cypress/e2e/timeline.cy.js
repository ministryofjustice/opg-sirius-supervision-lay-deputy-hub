describe("Timeline tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/lay/1");
        cy.get(':nth-child(3) > .moj-sub-navigation__link').click();
    });

    it("has header for section", () => {
        cy.get("h1").should("contain", "Timeline");
    });

    it("lists all types of timeline events", () => {
        cy.get('[data-cy="address-changed-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="fee-payer-added-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="fee-payer-removed-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="correspondence-document-published-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="death-notification-deleted-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="deputy-linked-to-order-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="deputy-order-details-changed-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="deputy-unlinked-from-order-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="person-additional-information-changed-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="person-contact-details-changed-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="person-death-notified-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="person-details-changed-event"] > .moj-timeline__header').should('be.visible');
        cy.get('[data-cy="person-status-changed-event"] > .moj-timeline__header').should('be.visible');
    });

    describe("Timeline event", () => {
        it("displays event title, byline, date and content", () => {
            cy.get('[data-cy="deputy-unlinked-from-order-event"] > .moj-timeline__header > .moj-timeline__title').should("contain", "Unlink");
            cy.get('[data-cy="deputy-unlinked-from-order-event"] > .moj-timeline__header > .moj-timeline__byline').should("contain", "by system admin (12345678)");
            cy.get('[data-cy="deputy-unlinked-from-order-event"] > .moj-timeline__date > time').should("contain", "13/10/2026");
            cy.get('[data-cy="deputy-unlinked-from-order-event"] > .moj-timeline__date > time').should("contain", "09:31:24");
            cy.get('[data-cy="deputy-unlinked-from-order-event"] > .moj-timeline__description > .govuk-list').should("contain", "Unlinked from case PFA");
        });
        it("displays changesets with set when old value is null", () => {
            cy.get('[data-cy="deputy-order-details-changed-event"] > .moj-timeline__description > .govuk-list').contains("set to Spouse");
        });
        it("displays changesets with changed from when old value is not null", () => {
            cy.get('[data-cy="deputy-order-details-changed-event"] > .moj-timeline__description > .govuk-list').contains("changed from Deceased to Error");
        });
    });

    describe("Amended information", () => {
        it("has filled in the timeline value for each changeset", () => {
            cy.get('[data-cy="deputy-order-details-changed-event"] > .moj-timeline__description').should("contain", "Type of deputy");
            cy.get('[data-cy="deputy-order-details-changed-event"] > .moj-timeline__description > .govuk-list').contains("Correspondent");
            cy.get('[data-cy="person-additional-information-changed-event"] > .moj-timeline__description').contains("Audio");
        });
        it("shows yes and no for boolean changesets", () => {
            cy.get('[data-cy="deputy-order-details-changed-event"] > .moj-timeline__description > .govuk-list').contains("changed from Yes to No");
            cy.get('[data-cy="person-additional-information-changed-event"] > .moj-timeline__description > .govuk-list').contains("set to Yes");
        });
    });

    describe("Functionality", () => {
        it("has hyperlinks to the client case", () => {
            cy.get('#hook-deputy-unlinked-from-order-timeline-event-client-id').should("be.visible", "7000-0000-1334");
            cy.get('#hook-deputy-unlinked-from-order-timeline-event-client-id').invoke("attr", "href").should("include", "/clients/47");
            cy.get('#hook-deputy-linked-to-order-timeline-event-client-id').should("be.visible", "7000-0000-1276");
            cy.get('#hook-deputy-linked-to-order-timeline-event-client-id').invoke("attr", "href").should("include", "/clients/47");
            cy.get('#hook-deputy-unlinked-from-order-timeline-event-client-id').should("be.visible", "7000-0000-1334");
            cy.get('#hook-deputy-unlinked-from-order-timeline-event-client-id').invoke("attr", "href").should("include", "/clients/47");
        });
    });
});
