describe("Manage deputy details tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    it("directs from main page and has header for section", () => {
        cy.visit("/supervision/deputies/lay/1");
        cy.get('#record-death').click();
        cy.get("h1").should("contain", "Record death of Mike Deputy");
    });

    it("has fillable content", () => {
        cy.visit("/supervision/deputies/lay/1/record-death");
        cy.get('input[name="proof-of-death"][value="true"]').check();
        cy.get('input[name="date-of-death"]').type("2020-01-01");
        cy.get('input[name="date-notified"]').type("2020-01-02");
        cy.get('input[name="date-proof-of-death-received"]').type("2020-01-03");
        cy.get('input[type="radio"][value="EXECUTOR"]').check();
        cy.get('input[type="radio"][value="PHONE"]').check();
    });

    describe("Autofill death record on form load", () => {
        it("fetches deputy death record into form and buttons direct to correct locations", () => {
            cy.visit("/supervision/deputies/lay/2/record-death");
            cy.get('#proof-of-death-true').should('be.checked');
            cy.get('[name="date-of-death"]').should('have.value', '2026-11-03');
            cy.get('[name="date-notified"]').should('have.value', '2026-11-05');
            cy.get('[name="date-proof-of-death-received"]').should('have.value', '2026-11-06');
            cy.get('#notified-by-other').should('be.checked');
            cy.get('#notification-method-phone').should('be.checked');
            cy.get('[name="notes"]').should('contain.text', 'Neighbour');

            cy.get('.govuk-button-group > .govuk-link').click();
            cy.get('.govuk-heading-l').should("contain.text", "Deputy details");

            cy.visit("/supervision/deputies/lay/3/record-death");
            cy.get('#proof-of-death-false').should('be.checked');
            cy.get('[name="date-of-death"]').should('have.value', '');
            cy.get('[name="date-notified"]').should('have.value', '2025-01-01');
            cy.get('[name="date-proof-of-death-received"]').should('have.value', '');
            cy.get('#notified-by-co-deputy').should('be.checked');
            cy.get('#notification-method-email').should('be.checked');
            cy.get('[name="notes"]').should('contain.text', 'Minimum death record content');
        });
    });
    describe("Deputy death record form validation messages", () => {
        it("Displays validation messages at header of page and above individual elements", () => {
            cy.setCookie("fail-route", "death-record");
            cy.visit("/supervision/deputies/lay/1/record-death");
            cy.get('.govuk-button').click();

            cy.get('.govuk-error-summary').should("be.visible");
            cy.get('.govuk-error-summary').should("contain", "Date notified cannot be a future date");
            cy.get('.govuk-error-summary').should("contain", "Date of death cannot be a future date");
            cy.get('.govuk-error-summary').should("contain", "The note must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "Notification method is required and can't be empty");
            cy.get('.govuk-error-summary').should("contain", "Notified by is required and can't be empty");
            cy.get('.govuk-error-summary').should("contain", "Proof of death received is required and can't be empty");

            cy.get('#proof-of-death-true').click();
            cy.get('#name-error-isEmpty').should("contain.text", "Proof of death received is required and can't be empty");
            cy.get('#f-dateOfDeath > #name-error-invalid-lte').should("contain", "Date of death cannot be a future date");
            cy.get('#f-dateNotified > #name-error-invalid-lte').should("contain", "Date notified cannot be a future date");
            cy.get('#f-dateDeathCertificateReceived > .govuk-form-group--error > #name-error-invalid-lte').should("contain", "Date proof of death received cannot be a future date");
            cy.get('#f-notifiedBy > #name-error-isEmpty').should("contain", "Notified by is required and can't be empty");
            cy.get('#f-notificationMethod > #name-error-isEmpty').should("contain", "Notification method is required and can't be empty");
            cy.get('#f-notes > .govuk-form-group > #name-error-stringLengthTooLong').should("contain", "The note must be 255 characters or fewer");
        })
    });
});
