describe("Clients tab", () => {
    beforeEach(() => {
        cy.setCookie("Other", "other");
        cy.setCookie("XSRF-TOKEN", "abcde");
    });

    it("should navigate to the Clients tab", () => {
        cy.visit("/supervision/deputies/lay/1");
        cy.get(".moj-sub-navigation__list").contains("Clients").click();
        cy.location("pathname").should("eq", "/supervision/deputies/lay/1/clients");
        cy.get("h1").should("contain", "Clients");
        cy.get(".moj-sub-navigation__link[aria-current='page']").should("have.text", "Clients");
    });

    describe("clients table", () => {
        beforeEach(() => {
            cy.visit("/supervision/deputies/lay/1/clients");
        });

        it("displays the seven column headings without filters", () => {
            const expected = [
                "Client",
                "Relationship",
                "Appointed",
                "Order made date",
                "Order type",
                "Order status",
                "Deputy status",
            ];

            cy.get(".govuk-table__head th")
                .should("have.length", expected.length)
                .each(($el, index) => {
                    cy.wrap($el).should("have.text", expected[index]);
                });
            cy.get("[data-module='filters']").should("not.exist");
        });

        it("lists clients with active, closed, open and duplicate orders", () => {
            const expected = [
                ["Agnes Burgundy", ["Sibling", "Sibling"], ["Joint and Severally", "Sole"], ["01/01/2016", "02/02/2016"], ["PFA", "HW"], ["Active", "Active"], ["Active", "Active"]],
                ["Louis Dauphin", ["Spouse", "Spouse"], ["Sole", "-"], ["01/01/2015", "01/01/2020"], ["PFA", "HW"], ["Active", "Closed"], ["Suspended", "Suspended"]],
                ["Something Here", "Parent", "Joint", "01/01/2016", "PFA", "Active", "Active"],
                ["Robert Robson", "Sibling", "Sole", "15/03/2024", "HW", "Closed", "Discharged"],
                ["Jane Smith", "Child", "Joint", "10/06/2024", "PFA", "Open", "Active"],
                ["Alex Taylor", "Spouse", "Joint and Severally", "01/09/2023", "PFA", "Duplicate", "Error"],
            ];

            cy.get(".govuk-table__body > .govuk-table__row")
                .should("have.length", expected.length)
                .each(($row, rowIndex) => {
                    cy.wrap($row).find("td").should("have.length", 7).each(($cell, columnIndex) => {
                        const value = expected[rowIndex][columnIndex];
                        if (columnIndex === 0) {
                            expect($cell.text().replace(/\s+/g, " ").trim()).to.contain(value);
                        } else {
                            const lines = $cell.children("div").toArray().map((line) => line.textContent.trim());
                            expect(lines).to.deep.equal(Array.isArray(value) ? value : [value]);
                        }
                    });
                });
        });

        it("links client names to the client hub with case references underneath", () => {
            const expected = [
                [71, "Agnes Burgundy", "93404674"],
                [74, "Louis Dauphin", "18227939"],
                [75, "Something Here", "93404699"],
                [76, "Robert Robson", "12323984"],
                [77, "Jane Smith", "12323985"],
                [78, "Alex Taylor", "12323986"],
            ];

            cy.get(".govuk-table__body > .govuk-table__row").each(($row, index) => {
                const [id, name, reference] = expected[index];
                cy.wrap($row).find("td").first().within(() => {
                    cy.get("a")
                        .should("have.length", 1)
                        .should("have.attr", "href")
                        .and("match", new RegExp(`/supervision/#/clients/${id}$`));
                    cy.get("a")
                        .invoke("text")
                        .then((text) => expect(text.replace(/\s+/g, " ").trim()).to.eq(name));
                    cy.get("br + span").should("have.length", 1).and("contain", reference);
                });
            });
        });

        it("aligns PFA details above HW details within the same client row", () => {
            cy.get(".govuk-table__body > .govuk-table__row").first().within(() => {
                cy.get("td").then(($cells) => {
                    for (let column = 1; column < $cells.length; column++) {
                        const lines = $cells.eq(column).children("div");
                        expect(lines).to.have.length(2);
                        expect(lines[0].getBoundingClientRect().top).to.eq(
                            $cells.eq(4).children("div")[0].getBoundingClientRect().top,
                        );
                        expect(lines[1].getBoundingClientRect().top).to.eq(
                            $cells.eq(4).children("div")[1].getBoundingClientRect().top,
                        );
                    }
                });
            });
        });

        it("does not wrap client names or appointment types", () => {
            cy.get(".govuk-table__body > .govuk-table__row").each(($row) => {
                cy.wrap($row).find("td").eq(0).should("have.css", "white-space", "nowrap");
                cy.wrap($row).find("td").eq(2).should("have.css", "white-space", "nowrap");
            });
        });

        it("keeps the deputy overview count limited to active clients", () => {
            cy.visit("/supervision/deputies/lay/1");
            cy.get("#overview").should("contain", "3").and("contain", "Active clients");
        });
    });
});
