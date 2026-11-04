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

        it("Redirects to the correct route", () => {
            cy.get('.govuk-button').click();
            cy.url().should("include", "?success=deputyDetails");
        })

    });
    describe("Deputy contact details form", () => {
        it("Displays validation messages", () => {
            cy.setCookie("fail-route", "contact-details");

            cy.get('.govuk-button').click();

            cy.get('.govuk-error-summary').should("be.visible");
            cy.get('.govuk-error-summary').should("contain", "The building or street must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "Address line 2 must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "Address line 3 must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The country must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The county must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The email must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The deputy first name is required and can't be empty");
            cy.get('.govuk-error-summary').should("contain", "The evening telephone number must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The deputy middle name must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The mobile number must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The deputy other name must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The postcode must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The deputy surname is required and can't be empty");
            cy.get('.govuk-error-summary').should("contain", "The town or city must be 255 characters or fewer");
            cy.get('.govuk-error-summary').should("contain", "The telephone number must be 255 characters or fewer");
        })
    });
});
