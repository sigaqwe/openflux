import UIKit

private let appTitle = "OpenFlux"
private let defaultURL = "https://docs.yandex.ru/"

@main
final class AppDelegate: UIResponder, UIApplicationDelegate {

    var window: UIWindow?

    func application(
        _ application: UIApplication,
        didFinishLaunchingWithOptions launchOptions:
        [UIApplication.LaunchOptionsKey: Any]? = nil
    ) -> Bool {

        let window = UIWindow(frame: UIScreen.main.bounds)
        window.rootViewController = OpenFluxViewController()
        window.makeKeyAndVisible()

        self.window = window

        return true
    }
}

final class OpenFluxViewController: UIViewController {

    private let urlField = UITextField()
    private let startButton = UIButton(type: .system)
    private let statusLabel = UILabel()
    private let proxyLabel = UILabel()

    private var started = false

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = .systemBackground

        title = appTitle

        setupUI()
    }

    private func setupUI() {

        let titleLabel = UILabel()
        titleLabel.text = "OpenFlux"
        titleLabel.font = .boldSystemFont(ofSize: 32)
        titleLabel.textAlignment = .center

        let descriptionLabel = UILabel()
        descriptionLabel.text =
            "OpenFlux client\n\nEnter your Yandex Docs transport URL:"
        descriptionLabel.numberOfLines = 0
        descriptionLabel.textAlignment = .center
        descriptionLabel.font = .systemFont(ofSize: 17)

        urlField.placeholder = "https://docs.yandex.ru/..."
        urlField.text = UserDefaults.standard.string(
            forKey: "openflux.url"
        ) ?? defaultURL
        urlField.borderStyle = .roundedRect
        urlField.autocapitalizationType = .none
        urlField.autocorrectionType = .no
        urlField.keyboardType = .URL
        urlField.returnKeyType = .done

        startButton.setTitle("Start OpenFlux", for: .normal)
        startButton.titleLabel?.font = .boldSystemFont(ofSize: 18)
        startButton.addTarget(
            self,
            action: #selector(startOpenFlux),
            for: .touchUpInside
        )

        statusLabel.text = "Status: stopped"
        statusLabel.textAlignment = .center
        statusLabel.numberOfLines = 0

        proxyLabel.text =
            "SOCKS5: 127.0.0.1:1080"
        proxyLabel.textAlignment = .center
        proxyLabel.textColor = .secondaryLabel
        proxyLabel.numberOfLines = 0

        let stack = UIStackView(arrangedSubviews: [
            titleLabel,
            descriptionLabel,
            urlField,
            startButton,
            statusLabel,
            proxyLabel
        ])

        stack.axis = .vertical
        stack.spacing = 18
        stack.translatesAutoresizingMaskIntoConstraints = false

        view.addSubview(stack)

        NSLayoutConstraint.activate([
            stack.leadingAnchor.constraint(
                equalTo: view.safeAreaLayoutGuide.leadingAnchor,
                constant: 24
            ),
            stack.trailingAnchor.constraint(
                equalTo: view.safeAreaLayoutGuide.trailingAnchor,
                constant: -24
            ),
            stack.centerYAnchor.constraint(
                equalTo: view.safeAreaLayoutGuide.centerYAnchor
            )
        ])
    }

    @objc private func startOpenFlux() {

        guard !started else {
            return
        }

        guard let text = urlField.text,
              let url = URL(string: text),
              url.scheme == "https" || url.scheme == "http"
        else {
            statusLabel.text = "Status: invalid URL"
            return
        }

        UserDefaults.standard.set(
            text,
            forKey: "openflux.url"
        )

        started = true

        urlField.isEnabled = false
        startButton.isEnabled = false

        statusLabel.text =
            "Status: starting...\n\nKeep OpenFlux open while using the proxy."

        DispatchQueue.global(qos: .userInitiated).async {

            text.withCString { pointer in
                RunMainClient(
                    UnsafeMutablePointer(
                        mutating: pointer
                    )
                )
            }

            DispatchQueue.main.async {
                self.statusLabel.text =
                    "Status: running\n\nSOCKS5: 127.0.0.1:1080"
            }
        }
    }
}
