describe("Manage additional information tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");

        cy.visit("/supervision/deputies/lay/1");
        cy.contains("Manage additional information").click();
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
            cy.get('[name="welsh"]').should("have.attr", "checked");

            cy.get('[name="additional-notes"]').should('contain', 'Deputy has additional notes');
        });

        it("saves the changes and redirects to the deputy details page", () => {
            cy.visit("/supervision/deputies/lay/2");
            cy.contains("Manage additional information").click();
            cy.get('[name="interpreter-required"]').clear().type("Updated interpreter requirements");
            cy.get('[name="audio"]').check();
            cy.get('[name="large-print"]').uncheck();
            cy.get('[name="hearing-impaired"]').uncheck();
            cy.get('[name="welsh"]').uncheck();
            cy.get('[name="additional-notes"]').clear().type("Updated additional information");

            cy.get('.govuk-button').click();
            cy.url().should("include", "?success=additionalInformation");

            cy.contains(".govuk-summary-list__key", "Interpreter requirements")
                .next(".govuk-summary-list__value")
                .should("contain.text", "Updated interpreter requirements");
            cy.contains(".govuk-summary-list__key", "Special correspondence requirements")
                .next(".govuk-summary-list__value")
                .should("contain.text", "Audio")
                .and("not.contain.text", "Large Print")
                .and("not.contain.text", "Hearing impaired")
                .and("not.contain.text", "Welsh");
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
