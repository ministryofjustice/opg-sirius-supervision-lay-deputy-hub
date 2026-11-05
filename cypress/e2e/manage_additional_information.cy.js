describe("Manage additional information tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/lay/1/manage-additional-information");
    });

    it("has header for section", () => {
        cy.get("h1").should("contain", "Manage additional information");
    });

    describe("Deputy additional information form", () => {
        it("fetches additional information into form", () => {
            cy.get('[name="interpreter-required"]').should('have.value', 'Interpreter is required for deputy');

            cy.get('[name="audio"]').should("have.attr", "checked");
            cy.get('[name="large-print"]').should("have.attr", "checked");
            cy.get('[name="hearing-impaired"]').should("have.attr", "checked");
            cy.get('[name="spelling-takes-care"]').should("have.attr", "checked");

            cy.get('[name="notes"]').should('contain', 'Deputy has additional notes');
        });

        it("Redirects to the correct route", () => {
            cy.get('.govuk-button').click();
            cy.url().should("include", "?success=additionalInformation");
        })

    });
    describe("Deputy additional information form", () => {
        it("Displays validation messages", () => {
            cy.setCookie("fail-route", "additional-information");

            cy.get('.govuk-button').click();

            cy.get('.govuk-error-summary').should("be.visible");
            cy.get('.govuk-error-summary').should("contain", "The interpreter requirements must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The additional information must be 1000 characters or fewer");
        })
    });
});
