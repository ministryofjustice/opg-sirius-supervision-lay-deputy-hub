describe("Manage deputy details tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
        cy.visit("/supervision/deputies/lay/1/manage-deputy-details");
    });

    it("has header for section", () => {
        cy.get("h1").should("contain", "Manage deputy contact details");
    });

    describe("Deputy contact details form", () => {
        it("fetches deputy details into form", () => {
            // cy.get("input #f-firstname").should("contain", "First name"); // Check labels?
            cy.get('[name="deputy-first-name"]').should('have.value', 'Mike');
            cy.get('[name="deputy-middle-name"]').should("have.value", "Jonathon");
            cy.get('[name="deputy-last-name"]').should("have.value", "Deputy");
            cy.get('[name="deputy-other-names"]').should("have.value", "Michael");

            cy.get('[name="deputy-date-of-birth"]').should("have.value", "1980-01-01");

            cy.get('[name="deputy-mobile-number"]').should("have.value", "0771 2345678");
            cy.get('[name="deputy-phone-number"]').should("have.value", "0115 876 5574");
            cy.get('[name="deputy-evening-number"]').should("have.value", "0115 2767825");
            cy.get('[name="deputy-email').should("have.value", "testemail@hotmail.co.uk");

            cy.get('[name="address-line-1"]').should("have.value", "Seax House");
            cy.get('[name="address-line-2"]').should("have.value", "19 Market Rd");
            cy.get('[name="address-line-3"]').should("have.value", "");
            cy.get('[name="town"]').should("have.value", "Chelmsford");
            cy.get('[name="county"]').should("have.value", "Essex");
            cy.get('[name="postcode"]').should("have.value", "CM1 1GG");
            cy.get('[name="country"]').should("have.value", "United Kingdom");

            cy.get('#airmail-required-false').should("have.attr", "checked");
        });

        it("Displays validation messages", () => {
            cy.setCookie("fail-route", "contact-details");

            // Required fields
            cy.get('[name="deputy-first-name"]').clear();
            cy.get('[name="deputy-last-name"]').clear();
            cy.get('[name="address-line-1"]').clear();
            cy.get('[name="town"]').clear();
            cy.get('[name="postcode"]').clear();

            // Over character limit
            cy.get('[name="deputy-other-names"]').clear().type('a'.repeat(256));

            cy.get('.govuk-button').click();

            cy.get('[data-module="govuk-error-summary"]').should("be.visible");
            cy.get('[data-module="govuk-error-summary"]').should("contain", "Enter a first name");
            cy.get('[data-module="govuk-error-summary"]').should("contain", "Enter a last name");
            cy.get('[data-module="govuk-error-summary"]').should("contain", "Enter an address");
            cy.get('[data-module="govuk-error-summary"]').should("contain", "Enter a town or city");
            cy.get('[data-module="govuk-error-summary"]').should("contain", "Enter a postcode");
            cy.get('[data-module="govuk-error-summary"]').should("contain", "Other names must be 255 characters or fewer");
        })
    });
});
