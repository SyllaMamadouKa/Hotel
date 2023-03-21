import gi
gi.require_version('Gtk', '3.0')
from gi.repository import Gtk


class Chambre:
    def __init__(self, id_chambre, num_chambre, prix):
        self.id_chambre = id_chambre
        self.num_chambre = num_chambre
        self.prix = prix

class ChambreWindow(Gtk.Window):
    def __init__(self):
        Gtk.Window.__init__(self, title="Ajouter une chambre")
        self.set_border_width(10)

        # Ajouter un conteneur pour les éléments de l'interface utilisateur
        box = Gtk.Box(orientation=Gtk.Orientation.VERTICAL, spacing=10)
        self.add(box)

        # Ajouter une étiquette pour le titre de la fenêtre
        title_label = Gtk.Label()
        title_label.set_markup("<big><b>Ajouter une chambre</b></big>")
        box.pack_start(title_label, False, False, 0)

        # Ajouter un champ de texte pour le nom
        nom_label = Gtk.Label("ID Chambre:")
        box.pack_start(nom_label, False, False, 0)
        self.nom_entry = Gtk.Entry()
        box.pack_start(self.nom_entry, False, False, 0)

        # Ajouter un champ de texte pour le numéro de chambre
        num_chambre_label = Gtk.Label("Numéro de chambre:")
        box.pack_start(num_chambre_label, False, False, 0)
        self.num_chambre_entry = Gtk.Entry()
        box.pack_start(self.num_chambre_entry, False, False, 0)

        # Ajouter un champ de texte pour le prix
        prix_label = Gtk.Label("Prix:")
        box.pack_start(prix_label, False, False, 0)
        self.prix_entry = Gtk.Entry()
        box.pack_start(self.prix_entry, False, False, 0)

        # Ajouter un bouton pour ajouter la chambre
        add_button = Gtk.Button.new_with_label("Ajouter")
        add_button.connect("clicked", self.on_add_button_clicked)
        box.pack_start(add_button, False, False, 0)

        # Afficher la fenêtre d'ajout de chambre
        self.show_all()

    def on_add_button_clicked(self, widget):
        # Récupérer les données du formulaire
        id_chambre = self.nom_entry.get_text()
        num_chambre = self.num
class Reservation:
    def __init__(self, nom, prenom, num_chambre, date_debut, date_fin, prix):
        self.nom = nom
        self.prenom = prenom
        self.num_chambre = num_chambre
        self.date_debut = date_debut
        self.date_fin = date_fin
        self.prix = prix

class ReservationWindow(Gtk.Window):
    def __init__(self):
        Gtk.Window.__init__(self, title="Ajouter une réservation")
        self.set_border_width(10)

        # Ajouter un conteneur pour les éléments de l'interface utilisateur
        box = Gtk.Box(orientation=Gtk.Orientation.VERTICAL, spacing=10)
        self.add(box)

        # Ajouter une étiquette pour le titre de la fenêtre
        title_label = Gtk.Label()
        title_label.set_markup("<big><b>Ajouter une réservation</b></big>")
        box.pack_start(title_label, False, False, 0)

        # Ajouter un champ de texte pour le nom
        nom_label = Gtk.Label("Nom:")
        box.pack_start(nom_label, False, False, 0)
        self.nom_entry = Gtk.Entry()
        box.pack_start(self.nom_entry, False, False, 0)

        # Ajouter un champ de texte pour le prénom
        prenom_label = Gtk.Label("Prénom:")
        box.pack_start(prenom_label, False, False, 0)
        self.prenom_entry = Gtk.Entry()
        box.pack_start(self.prenom_entry, False, False, 0)

        # Ajouter un champ de texte pour le numéro de chambre
        num_chambre_label = Gtk.Label("Numéro de chambre:")
        box.pack_start(num_chambre_label, False, False, 0)
        self.num_chambre_entry = Gtk.Entry()
        box.pack_start(self.num_chambre_entry, False, False, 0)

        # Ajouter un champ de texte pour la date de début
        date_debut_label = Gtk.Label("Date de début:")
        box.pack_start(date_debut_label, False, False, 0)
        self.date_debut_entry = Gtk.Entry()
        box.pack_start(self.date_debut_entry, False, False, 0)

        # Ajouter un champ de texte pour la date de fin
        date_fin_label = Gtk.Label("Date de fin:")
        box.pack_start(date_fin_label, False, False, 0)
        self.date_fin_entry = Gtk.Entry()
        box.pack_start(self.date_fin_entry, False, False, 0)

        # Ajouter un champ de texte pour le prix
        prix_label = Gtk.Label("Prix:")
        box.pack_start(prix_label, False, False, 0)
        self.prix_entry = Gtk.Entry()
        box.pack_start(self.prix_entry, False, False, 0)

        # Ajouter un bouton pour ajouter la réservation
        add_button = Gtk.Button.new_with_label("Ajouter")
        add_button.connect("clicked", self.on_add_button_clicked)
        box.pack_start(add_button, False, False, 0)

        # Afficher la fenêtre d'ajout de réservation
        self.show_all()

    def on_add_button_clicked(self, widget):
        # Récupérer les données du formulaire
        nom = self.nom_entry.get_text()
        prenom = self.prenom_entry.get_text()
        num_chambre = self.num

class MainWindow(Gtk.Window):
    def __init__(self):
        Gtk.Window.__init__(self, title="Gestion d'un hôtel")
        self.set_border_width(10)

        # Ajouter un conteneur pour les éléments de l'interface utilisateur
        box = Gtk.Box(orientation=Gtk.Orientation.VERTICAL, spacing=10)
        self.add(box)

        # Ajouter une étiquette pour le titre de l'application
        title_label = Gtk.Label()
        title_label.set_markup("<big><b>Gestion d'un hôtel</b></big>")
        box.pack_start(title_label, False, False, 0)

        # Ajouter un bouton pour la gestion des réservations
        reservations_button = Gtk.Button.new_with_label("Gérer les réservations")
        reservations_button.connect("clicked", self.on_reservations_button_clicked)
        box.pack_start(reservations_button, False, False, 0)

        # Ajouter un bouton pour la gestion des chambres
        rooms_button = Gtk.Button.new_with_label("Gérer les chambres")
        rooms_button.connect("clicked", self.on_rooms_button_clicked)
        box.pack_start(rooms_button, False, False, 0)

        # Ajouter un bouton pour la gestion des factures
        invoices_button = Gtk.Button.new_with_label("Gérer les factures")
        invoices_button.connect("clicked", self.on_invoices_button_clicked)
        box.pack_start(invoices_button, False, False, 0)

        # Afficher la fenêtre principale
        self.show_all()

    def on_reservations_button_clicked(self, widget):
        # Ouvrir la fenêtre de réservation
        reservation_window = ReservationWindow()
        reservation_window.show_all()
        pass

    def on_rooms_button_clicked(self, widget):
        # Gérer les chambres
        # Ouvrir la fenêtre de réservation
        chambre_window = ChambreWindow()
        chambre_window.show_all()
        pass
        pass

    def on_invoices_button_clicked(self, widget):
        # Gérer les factures
        pass

if __name__ == "__main__":
    win = MainWindow()
    win.connect("delete-event", Gtk.main_quit)
    Gtk.main()
